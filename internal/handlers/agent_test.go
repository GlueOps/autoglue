package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func agentReq(t *testing.T, method, target string, body any, agent *models.Agent, params map[string]string) *http.Request {
	t.Helper()
	var rdr *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	r := httptest.NewRequest(method, target, rdr)
	ctx := r.Context()
	if agent != nil {
		ctx = httpmiddleware.WithAgent(ctx, agent)
	}
	rc := chi.NewRouteContext()
	for k, v := range params {
		rc.URLParams.Add(k, v)
	}
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rc)
	return r.WithContext(ctx)
}

func seedAgentWorld(t *testing.T, db *gorm.DB) (org models.Organization, cluster models.Cluster, server models.Server) {
	t.Helper()
	org = models.Organization{ID: uuid.New(), Name: "o-" + uuid.NewString()[:8]}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("org: %v", err)
	}
	key := models.SshKey{
		Name: "k", PublicKey: "p", EncryptedPrivateKey: "e",
		PrivateIV: "i", PrivateTag: "t", Fingerprint: "f",
	}
	key.OrganizationID = org.ID
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("sshkey: %v", err)
	}
	server = models.Server{
		OrganizationID: org.ID, Hostname: "bastion", SSHUser: "deploy",
		Status: "ready", SshKeyID: key.ID,
	}
	if err := db.Create(&server).Error; err != nil {
		t.Fatalf("server: %v", err)
	}
	cluster = models.Cluster{
		OrganizationID: org.ID, Name: "c1",
		DockerImage: "ghcr.io/glueops/kube", DockerTag: "v1",
		BastionServerID: &server.ID,
	}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("cluster: %v", err)
	}
	return org, cluster, server
}

func mintTicket(t *testing.T, db *gorm.DB, org models.Organization, c models.Cluster, s models.Server) string {
	t.Helper()
	plain, err := auth.RandomB64URL(32)
	if err != nil {
		t.Fatal(err)
	}
	tk := models.AgentEnrollmentTicket{
		OrganizationID: org.ID, ClusterID: c.ID, ServerID: s.ID,
		TicketHash: auth.SHA256Hex(plain), Prefix: plain[:12],
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := db.Create(&tk).Error; err != nil {
		t.Fatalf("ticket: %v", err)
	}
	return plain
}

func TestAgentPlaneFullCycle(t *testing.T) {
	db := pgtest.DB(t)
	org, cluster, server := seedAgentWorld(t, db)

	// --- enrol
	plain := mintTicket(t, db, org, cluster, server)
	rr := httptest.NewRecorder()
	EnrollAgent(db)(rr, agentReq(t, http.MethodPost, "/agent/enroll",
		agentEnrollRequest{Ticket: plain, Version: "0.1.0"}, nil, nil))
	if rr.Code != http.StatusCreated {
		t.Fatalf("enroll = %d %s", rr.Code, rr.Body.String())
	}
	var enrolled agentEnrollResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &enrolled); err != nil {
		t.Fatal(err)
	}
	if auth.ValidateAgentKeyPair(enrolled.AgentID.String(), enrolled.Key, enrolled.Secret, db) == nil {
		t.Fatal("minted credential does not validate")
	}

	// replay of the same ticket must be a flat 401
	rr = httptest.NewRecorder()
	EnrollAgent(db)(rr, agentReq(t, http.MethodPost, "/agent/enroll",
		agentEnrollRequest{Ticket: plain}, nil, nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("ticket replay = %d, want 401", rr.Code)
	}

	var agent models.Agent
	if err := db.First(&agent, "id = ?", enrolled.AgentID).Error; err != nil {
		t.Fatal(err)
	}

	// --- sync with nothing published
	rr = httptest.NewRecorder()
	AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?generation=0", nil, &agent, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("sync (unpublished) = %d, want 204", rr.Code)
	}

	// publish generation 3
	if err := db.Create(&models.ClusterDesiredState{
		OrganizationID: org.ID, ClusterID: cluster.ID, Generation: 3, PublishedAt: time.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.DesiredResource{
		OrganizationID: org.ID, ClusterID: cluster.ID, Generation: 3,
		ResourceType: "kubeconfig", ResourceID: "default", Phase: 1,
		Spec: []byte(`{"a":1}`), DependsOn: []byte(`["x"]`), SpecHash: "h", Required: true,
	}).Error; err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?generation=2", nil, &agent, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("sync = %d %s", rr.Code, rr.Body.String())
	}
	var snap agentSnapshot
	if err := json.Unmarshal(rr.Body.Bytes(), &snap); err != nil {
		t.Fatal(err)
	}
	if snap.Generation != 3 || !snap.FullSnapshot || len(snap.Resources) != 1 {
		t.Fatalf("snapshot = %+v", snap)
	}
	if snap.Resources[0].Spec != `{"a": 1}` && snap.Resources[0].Spec != `{"a":1}` {
		t.Fatalf("spec = %q, want a JSON string", snap.Resources[0].Spec)
	}

	rr = httptest.NewRecorder()
	AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?generation=3", nil, &agent, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("sync (current) = %d, want 204", rr.Code)
	}

	// --- assignment when idle
	rr = httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
	if rr.Code != http.StatusOK || bytes.TrimSpace(rr.Body.Bytes())[0] != 'n' {
		t.Fatalf("idle assignment = %d %q", rr.Code, rr.Body.String())
	}

	// --- a run appears
	run := models.ClusterRun{
		OrganizationID: org.ID, ClusterID: cluster.ID,
		Action: "bootstrap", Status: models.ClusterRunStatusQueued,
		// Agent-owned on purpose: a run left as the default river executor is
		// one the control plane is executing over SSH, and the agent must not
		// claim it.
		Executor: models.ClusterRunExecutorAgent,
	}
	if err := db.Create(&run).Error; err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("assignment = %d %s", rr.Code, rr.Body.String())
	}
	var as agentAssignmentResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &as); err != nil {
		t.Fatal(err)
	}
	if as.Command != "bootstrap" || as.State != models.AgentTaskStateAssigned {
		t.Fatalf("assignment = %+v", as)
	}
	var ta agentTaskArgs
	if err := json.Unmarshal([]byte(as.Args), &ta); err != nil {
		t.Fatalf("args not a JSON string: %v (%q)", err, as.Args)
	}
	if ta.Image != "ghcr.io/glueops/kube" || ta.RunID != run.ID.String() || len(ta.Mounts) != 2 {
		t.Fatalf("task args = %+v", ta)
	}
	if ta.Mounts[0] != "/home/deploy/.ssh:/root/.ssh" {
		t.Fatalf("mounts = %v", ta.Mounts)
	}

	// the run flipped to running
	var afterRun models.ClusterRun
	db.First(&afterRun, "id = ?", run.ID)
	if afterRun.Status != models.ClusterRunStatusRunning {
		t.Fatalf("run status = %s", afterRun.Status)
	}

	// --- assignment restates the SAME task, and does not mint a second one
	for i := 0; i < 3; i++ {
		rr = httptest.NewRecorder()
		AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
		var again agentAssignmentResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &again)
		if again.TaskID != as.TaskID || again.State != models.AgentTaskStateAssigned {
			t.Fatalf("restate %d = %+v, want %s/assigned", i, again, as.TaskID)
		}
	}
	// Scoped to this cluster: pgtest hands every test in the package the same
	// database with no truncation between them, so an unscoped count here would
	// pass or fail depending on which tests ran first.
	var nTasks int64
	db.Model(&models.AgentTask{}).Where("cluster_id = ?", cluster.ID).Count(&nTasks)
	if nTasks != 1 {
		t.Fatalf("task count = %d, want 1", nTasks)
	}

	// a second queued run must NOT produce a second in-flight task
	run2 := models.ClusterRun{
		OrganizationID: org.ID, ClusterID: cluster.ID,
		Action: "upgrade", Status: models.ClusterRunStatusQueued,
		Executor: models.ClusterRunExecutorAgent,
	}
	db.Create(&run2)
	rr = httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
	var stillFirst agentAssignmentResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &stillFirst)
	if stillFirst.TaskID != as.TaskID {
		t.Fatalf("second run stole the slot: %+v", stillFirst)
	}

	params := map[string]string{"taskID": as.TaskID}

	// --- start, then a redelivered start, then a conflicting one
	rr = httptest.NewRecorder()
	AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
		agentTaskStartRequest{ContainerID: "deadbeef"}, &agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("start = %d %s", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
		agentTaskStartRequest{ContainerID: "deadbeef"}, &agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("start redelivery = %d, want 204", rr.Code)
	}
	rr = httptest.NewRecorder()
	AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
		agentTaskStartRequest{ContainerID: "other"}, &agent, params))
	if rr.Code != http.StatusConflict {
		t.Fatalf("start conflict = %d, want 409", rr.Code)
	}

	// assignment now reports started
	rr = httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
	var started agentAssignmentResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &started)
	if started.State != models.AgentTaskStateStarted {
		t.Fatalf("state = %s, want started", started.State)
	}

	// --- logs: 1..4, then redelivery of 1..7
	rr = httptest.NewRecorder()
	AgentTaskLogs(db)(rr, agentReq(t, http.MethodPost, "/logs", agentTaskLogsRequest{
		Chunks: []agentLogChunk{
			{Seq: 1, Stream: "stdout", Chunk: "a"},
			{Seq: 2, Stream: "stderr", Chunk: "b"},
			{Seq: 3, Stream: "weird", Chunk: "c"},
			{Seq: 4, Stream: "stdout", Chunk: "d"},
		},
	}, &agent, params))
	if rr.Code != http.StatusOK {
		t.Fatalf("logs = %d %s", rr.Code, rr.Body.String())
	}
	var lres agentTaskLogsResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &lres)
	if lres.AcceptedThroughSeq != 4 {
		t.Fatalf("accepted = %d, want 4", lres.AcceptedThroughSeq)
	}

	rr = httptest.NewRecorder()
	AgentTaskLogs(db)(rr, agentReq(t, http.MethodPost, "/logs", agentTaskLogsRequest{
		Chunks: []agentLogChunk{
			{Seq: 1, Stream: "stdout", Chunk: "a"},
			{Seq: 2, Stream: "stderr", Chunk: "b"},
			{Seq: 3, Stream: "weird", Chunk: "c"},
			{Seq: 4, Stream: "stdout", Chunk: "d"},
			{Seq: 5, Stream: "stdout", Chunk: "e"},
			{Seq: 6, Stream: "stdout", Chunk: "f"},
			{Seq: 7, Stream: "stdout", Chunk: "g"},
		},
	}, &agent, params))
	_ = json.Unmarshal(rr.Body.Bytes(), &lres)
	if lres.AcceptedThroughSeq != 7 {
		t.Fatalf("accepted = %d, want 7", lres.AcceptedThroughSeq)
	}

	var logs []models.JobLog
	db.Where("subject_type = ? AND subject_id = ?", models.JobLogSubjectClusterRun, run.ID).
		Order("id ASC").Find(&logs)
	if len(logs) != 7 {
		t.Fatalf("job_logs = %d, want 7 (no duplicates)", len(logs))
	}
	if logs[0].OrganizationID != org.ID || logs[1].Stream != "stderr" || logs[2].Stream != "stdout" {
		t.Fatalf("log rows = %+v", logs[:3])
	}

	// the existing reader must see them unchanged
	items, _, err := readJobLogs(db, org.ID, models.JobLogSubjectClusterRun, run.ID, 0, 100)
	if err != nil || len(items) != 7 {
		t.Fatalf("readJobLogs = %d %v", len(items), err)
	}

	// --- result: succeeded
	rr = httptest.NewRecorder()
	code := 0
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{State: "succeeded", ExitCode: &code}, &agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("result = %d %s", rr.Code, rr.Body.String())
	}
	db.First(&afterRun, "id = ?", run.ID)
	if afterRun.Status != "succeeded" || afterRun.FinishedAt.IsZero() {
		t.Fatalf("run after success = %+v", afterRun)
	}

	// redelivery is 204, contradiction is 409
	rr = httptest.NewRecorder()
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{State: "succeeded"}, &agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("result redelivery = %d, want 204", rr.Code)
	}
	rr = httptest.NewRecorder()
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{State: "failed"}, &agent, params))
	if rr.Code != http.StatusConflict {
		t.Fatalf("result contradiction = %d, want 409", rr.Code)
	}
	rr = httptest.NewRecorder()
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{State: "banana"}, &agent, params))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("bad state = %d, want 400", rr.Code)
	}

	// --- the slot is free, so run2 is now assignable
	rr = httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, &agent, nil))
	var second agentAssignmentResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &second)
	if second.Command != "upgrade" || second.TaskID == as.TaskID {
		t.Fatalf("second assignment = %+v", second)
	}

	// --- reconcile report
	rr = httptest.NewRecorder()
	AgentReconcileReport(db)(rr, agentReq(t, http.MethodPost, "/reconcile-report",
		agentReconcileReportRequest{
			CurrentGeneration: 3, AppliedGeneration: 3, Healthy: true,
			Resources: []agentResourceStatus{
				{ResourceType: "kubeconfig", ResourceID: "default", DesiredGeneration: 3, State: "applied"},
			},
		}, &agent, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("reconcile = %d %s", rr.Code, rr.Body.String())
	}
	// upsert, not duplicate
	rr = httptest.NewRecorder()
	AgentReconcileReport(db)(rr, agentReq(t, http.MethodPost, "/reconcile-report",
		agentReconcileReportRequest{
			CurrentGeneration: 3, AppliedGeneration: 0, Healthy: true,
			Resources: []agentResourceStatus{
				{ResourceType: "kubeconfig", ResourceID: "default", DesiredGeneration: 3, State: "failed", LastError: "boom"},
			},
		}, &agent, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("reconcile 2 = %d %s", rr.Code, rr.Body.String())
	}
	var statuses []models.AgentReconcileStatus
	db.Where("agent_id = ?", agent.ID).Find(&statuses)
	if len(statuses) != 1 || statuses[0].State != "failed" || statuses[0].LastError != "boom" {
		t.Fatalf("statuses = %+v", statuses)
	}
	var reloaded models.Agent
	db.First(&reloaded, "id = ?", agent.ID)
	if reloaded.AppliedGeneration != 3 {
		t.Fatalf("applied_generation walked backwards: %d", reloaded.AppliedGeneration)
	}

	// --- re-enrolment dead-letters the in-flight task and fails its run
	plain2 := mintTicket(t, db, org, cluster, server)
	rr = httptest.NewRecorder()
	EnrollAgent(db)(rr, agentReq(t, http.MethodPost, "/agent/enroll",
		agentEnrollRequest{Ticket: plain2}, nil, nil))
	if rr.Code != http.StatusCreated {
		t.Fatalf("re-enroll = %d %s", rr.Code, rr.Body.String())
	}
	var dl models.AgentTask
	db.First(&dl, "id = ?", second.TaskID)
	if dl.State != models.AgentTaskStateDeadLettered ||
		dl.DeadLetterReason != models.AgentTaskDeadLetterReenrolled {
		t.Fatalf("dead letter = %+v", dl)
	}
	var run2After models.ClusterRun
	db.First(&run2After, "id = ?", run2.ID)
	if run2After.Status != models.ClusterRunStatusFailed ||
		run2After.Error == "" || run2After.FinishedAt.IsZero() {
		t.Fatalf("run2 = %+v", run2After)
	}
	t.Logf("run2 error = %q", run2After.Error)

	// old credential is dead, new one lives, and only one is active
	if auth.ValidateAgentKeyPair(enrolled.AgentID.String(), enrolled.Key, enrolled.Secret, db) != nil {
		t.Fatal("revoked agent still authenticates")
	}
	var active int64
	db.Model(&models.Agent{}).Where("cluster_id = ? AND status = ?", cluster.ID, models.AgentStatusActive).Count(&active)
	if active != 1 {
		t.Fatalf("active agents = %d, want 1", active)
	}
}

// ---------------------------------------------------------------------------
// isolation, slot discipline and idempotency
//
// Everything below tests something that fails *quietly*. The full-cycle test
// above proves the happy path works; these prove the ways it can go wrong
// without anyone noticing until a cluster has been acted on by the wrong
// bastion, or a fleet has deadlocked, or a log tail has gone to nonsense.
// ---------------------------------------------------------------------------

// seedEnrolledAgent writes an active agent straight to the table. The enrolment
// handler is exercised in the full-cycle test; these tests care about what an
// authenticated agent may then reach, so they skip the ticket dance.
func seedEnrolledAgent(t *testing.T, db *gorm.DB, org models.Organization, cluster models.Cluster, server models.Server) *models.Agent {
	t.Helper()
	cred, err := auth.MintAgentCredential()
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	a := models.Agent{
		ID:             cred.ID,
		OrganizationID: org.ID,
		ClusterID:      cluster.ID,
		ServerID:       server.ID,
		KeyHash:        cred.KeyHash,
		SecretHash:     cred.SecretHash,
		Prefix:         cred.Prefix,
		Status:         models.AgentStatusActive,
		EnrolledAt:     time.Now().UTC(),
	}
	if err := db.Create(&a).Error; err != nil {
		t.Fatalf("agent: %v", err)
	}
	return &a
}

// seedSecondCluster adds another cluster to an existing org, sharing its
// bastion. Sharper than a second org for scoping tests: an authorization bug
// that filters on the organization alone still passes a cross-org check.
func seedSecondCluster(t *testing.T, db *gorm.DB, org models.Organization, server models.Server) models.Cluster {
	t.Helper()
	c := models.Cluster{
		OrganizationID: org.ID, Name: "c2-" + uuid.NewString()[:8],
		DockerImage: "ghcr.io/glueops/kube", DockerTag: "v1",
		BastionServerID: &server.ID,
	}
	if err := db.Create(&c).Error; err != nil {
		t.Fatalf("cluster: %v", err)
	}
	return c
}

func seedQueuedRun(t *testing.T, db *gorm.DB, org models.Organization, cluster models.Cluster, action string) models.ClusterRun {
	t.Helper()
	run := models.ClusterRun{
		OrganizationID: org.ID, ClusterID: cluster.ID,
		Action: action, Status: models.ClusterRunStatusQueued,
		Executor: models.ClusterRunExecutorAgent,
	}
	if err := db.Create(&run).Error; err != nil {
		t.Fatalf("run: %v", err)
	}
	return run
}

func assignOnce(t *testing.T, db *gorm.DB, agent *models.Agent) (agentAssignmentResponse, bool) {
	t.Helper()
	rr := httptest.NewRecorder()
	AgentAssignment(db)(rr, agentReq(t, http.MethodGet, "/agent/assignment", nil, agent, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("assignment = %d %s", rr.Code, rr.Body.String())
	}
	body := bytes.TrimSpace(rr.Body.Bytes())
	if string(body) == "null" {
		return agentAssignmentResponse{}, false
	}
	var as agentAssignmentResponse
	if err := json.Unmarshal(body, &as); err != nil {
		t.Fatalf("assignment body %q: %v", body, err)
	}
	return as, true
}

// TestAgentCannotTouchAnotherClustersTask.
//
// What breaks in production if this fails: one tenant's bastion writes the
// outcome of another tenant's cluster run. The three task endpoints take a task
// id straight out of the URL, and the only thing standing between a guessed or
// leaked uuid and someone else's run is the `agent_id = ?` predicate that is
// ANDed onto the lookup. Drop it and every failure mode is silent from the
// caller's side: the victim's run is marked failed, or succeeded, by a machine
// that never touched it, and its log tail fills with another cluster's output.
func TestAgentCannotTouchAnotherClustersTask(t *testing.T) {
	db := pgtest.DB(t)

	orgA, clusterA, serverA := seedAgentWorld(t, db)
	agentA := seedEnrolledAgent(t, db, orgA, clusterA, serverA)

	orgB, clusterB, serverB := seedAgentWorld(t, db)
	agentB := seedEnrolledAgent(t, db, orgB, clusterB, serverB)
	runB := seedQueuedRun(t, db, orgB, clusterB, "bootstrap")
	taskB := models.AgentTask{
		OrganizationID: orgB.ID, ClusterID: clusterB.ID, AgentID: agentB.ID,
		RunID: runB.ID, Kind: models.AgentTaskKindContainer, Command: "bootstrap",
		State: models.AgentTaskStateAssigned, AssignedAt: time.Now().UTC(),
	}
	if err := db.Create(&taskB).Error; err != nil {
		t.Fatal(err)
	}

	victim := map[string]string{"taskID": taskB.ID.String()}

	for name, call := range map[string]func(*httptest.ResponseRecorder){
		"start": func(rr *httptest.ResponseRecorder) {
			AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
				agentTaskStartRequest{ContainerID: "stolen"}, agentA, victim))
		},
		"logs": func(rr *httptest.ResponseRecorder) {
			AgentTaskLogs(db)(rr, agentReq(t, http.MethodPost, "/logs", agentTaskLogsRequest{
				Chunks: []agentLogChunk{{Seq: 1, Stream: "stdout", Chunk: "leaked"}},
			}, agentA, victim))
		},
		"result": func(rr *httptest.ResponseRecorder) {
			AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
				agentTaskResultRequest{State: models.AgentTaskStateFailed, Error: "not mine"}, agentA, victim))
		},
	} {
		rr := httptest.NewRecorder()
		call(rr)
		// 404 rather than 403 on purpose: as far as this principal is
		// concerned that task does not exist, and saying "forbidden" would
		// confirm the id names something real.
		if rr.Code != http.StatusNotFound {
			t.Errorf("%s on another cluster's task = %d, want 404 (%s)", name, rr.Code, rr.Body.String())
		}
	}

	var after models.AgentTask
	if err := db.First(&after, "id = ?", taskB.ID).Error; err != nil {
		t.Fatal(err)
	}
	if after.State != models.AgentTaskStateAssigned || after.ContainerID != "" || after.LogSeq != 0 {
		t.Fatalf("another agent moved the task: %+v", after)
	}
	var afterRun models.ClusterRun
	if err := db.First(&afterRun, "id = ?", runB.ID).Error; err != nil {
		t.Fatal(err)
	}
	if afterRun.Status != models.ClusterRunStatusQueued {
		t.Fatalf("another agent settled the run: %q", afterRun.Status)
	}
	var nLogs int64
	db.Model(&models.JobLog{}).
		Where("subject_type = ? AND subject_id = ?", models.JobLogSubjectClusterRun, runB.ID).
		Count(&nLogs)
	if nLogs != 0 {
		t.Fatalf("%d log rows written into another cluster's run", nLogs)
	}
}

// TestAgentAssignmentIsScopedToItsOwnCluster.
//
// What breaks in production if this fails: an agent runs a make target against
// the wrong cluster. Assignment is derived, not requested — the agent asks
// "what should I be doing" and does whatever comes back — so a scope bug here
// is not a leak of data, it is a bastion applying one cluster's bootstrap to
// another and the control plane recording it as a success.
//
// The second cluster deliberately shares the organization. A derivation that
// filtered on organization_id alone would pass a cross-tenant test and still
// hand this agent its neighbour's work.
func TestAgentAssignmentIsScopedToItsOwnCluster(t *testing.T) {
	db := pgtest.DB(t)

	org, clusterA, server := seedAgentWorld(t, db)
	agentA := seedEnrolledAgent(t, db, org, clusterA, server)

	clusterB := seedSecondCluster(t, db, org, server)
	runB := seedQueuedRun(t, db, org, clusterB, "upgrade")

	otherOrg, otherCluster, otherServer := seedAgentWorld(t, db)
	_ = seedEnrolledAgent(t, db, otherOrg, otherCluster, otherServer)
	otherRun := seedQueuedRun(t, db, otherOrg, otherCluster, "bootstrap")

	if as, ok := assignOnce(t, db, agentA); ok {
		t.Fatalf("agent was handed work it does not own: %+v", as)
	}

	for name, id := range map[string]uuid.UUID{
		"same-org other cluster": runB.ID,
		"other org":              otherRun.ID,
	} {
		var n int64
		db.Model(&models.AgentTask{}).Where("run_id = ?", id).Count(&n)
		if n != 0 {
			t.Errorf("%s: %d tasks minted for a run this agent cannot see", name, n)
		}
	}

	// Sanity anchor: the same agent does pick up work on its own cluster, so a
	// green result above is scoping rather than a derivation that never fires.
	own := seedQueuedRun(t, db, org, clusterA, "bootstrap")
	as, ok := assignOnce(t, db, agentA)
	if !ok || as.Command != "bootstrap" {
		t.Fatalf("agent got no work on its own cluster: %+v", as)
	}
	var task models.AgentTask
	if err := db.First(&task, "id = ?", as.TaskID).Error; err != nil {
		t.Fatal(err)
	}
	if task.RunID != own.ID || task.ClusterID != clusterA.ID {
		t.Fatalf("task points at the wrong run/cluster: %+v", task)
	}
}

// TestAgentSyncIsScopedAndDeliversExactlyOneGeneration.
//
// What breaks in production if this fails: the config plane hands a bastion the
// wrong desired state, which is worse than handing it none. Two ways, both
// invisible in a response that looks perfectly well-formed:
//
//   - Dropping the cluster predicate delivers another cluster's credentials and
//     shape to this bastion, which will then reconcile its disk to them.
//   - Dropping the generation predicate merges every generation ever published
//     into one snapshot. Resources removed in a later generation come back, and
//     because the snapshot is authoritative the agent will never delete
//     anything again.
func TestAgentSyncIsScopedAndDeliversExactlyOneGeneration(t *testing.T) {
	db := pgtest.DB(t)

	org, clusterA, server := seedAgentWorld(t, db)
	agentA := seedEnrolledAgent(t, db, org, clusterA, server)
	clusterB := seedSecondCluster(t, db, org, server)

	publish := func(c models.Cluster, gen int64, ids ...string) {
		t.Helper()
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "cluster_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"generation", "published_at"}),
		}).Create(&models.ClusterDesiredState{
			OrganizationID: org.ID, ClusterID: c.ID, Generation: gen, PublishedAt: time.Now(),
		}).Error; err != nil {
			t.Fatalf("publish: %v", err)
		}
		for _, id := range ids {
			if err := db.Create(&models.DesiredResource{
				OrganizationID: org.ID, ClusterID: c.ID, Generation: gen,
				ResourceType: "kubeconfig", ResourceID: id, Phase: 1,
				Spec: []byte(`{}`), SpecHash: "h", Required: true,
			}).Error; err != nil {
				t.Fatalf("resource: %v", err)
			}
		}
	}

	// Cluster B is ahead and has resources cluster A never had.
	publish(clusterB, 9, "b-only-1", "b-only-2")

	// Nothing published for A yet: 204, not an empty snapshot. An empty
	// resource list is a complete desired state meaning "delete everything you
	// have", so a 200 here would wipe a working bastion.
	rr := httptest.NewRecorder()
	AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?generation=0", nil, agentA, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("sync before publish = %d %s", rr.Code, rr.Body.String())
	}

	publish(clusterA, 1, "gone-in-2", "kept")
	publish(clusterA, 2, "kept")

	rr = httptest.NewRecorder()
	AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?generation=1", nil, agentA, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("sync = %d %s", rr.Code, rr.Body.String())
	}
	var snap agentSnapshot
	if err := json.Unmarshal(rr.Body.Bytes(), &snap); err != nil {
		t.Fatal(err)
	}
	if snap.Generation != 2 {
		t.Fatalf("generation = %d, want 2", snap.Generation)
	}
	if len(snap.Resources) != 1 || snap.Resources[0].ResourceID != "kept" {
		t.Fatalf("snapshot carries stale or foreign resources: %+v", snap.Resources)
	}
}

// TestAgentSyncReturns204WhenTheAgentIsCurrentOrAhead.
//
// What breaks in production if this fails: nothing user-visible, which is the
// problem. Every bastion re-downloads and re-applies a full desired state every
// poll interval forever, and the only symptom is load — on the API, on the
// database, and on the reconcile loop that now churns instead of idling. A
// comparison written as `!=` or `>` instead of `<=` also breaks the ahead case,
// which happens for real whenever a generation is rolled back.
func TestAgentSyncReturns204WhenTheAgentIsCurrentOrAhead(t *testing.T) {
	db := pgtest.DB(t)

	org, cluster, server := seedAgentWorld(t, db)
	agent := seedEnrolledAgent(t, db, org, cluster, server)

	if err := db.Create(&models.ClusterDesiredState{
		OrganizationID: org.ID, ClusterID: cluster.ID, Generation: 5, PublishedAt: time.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.DesiredResource{
		OrganizationID: org.ID, ClusterID: cluster.ID, Generation: 5,
		ResourceType: "kubeconfig", ResourceID: "default", Phase: 1,
		Spec: []byte(`{}`), SpecHash: "h", Required: true,
	}).Error; err != nil {
		t.Fatal(err)
	}

	for name, q := range map[string]string{
		"exactly current": "generation=5",
		"ahead":           "generation=6",
	} {
		rr := httptest.NewRecorder()
		AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?"+q, nil, agent, nil))
		if rr.Code != http.StatusNoContent {
			t.Errorf("%s = %d, want 204 — every poll re-delivers the whole snapshot", name, rr.Code)
		}
		if rr.Body.Len() != 0 {
			t.Errorf("%s returned a body with 204: %q", name, rr.Body.String())
		}
	}

	// Behind, or unable to state a generation at all, must still be served.
	for name, q := range map[string]string{
		"behind":      "generation=4",
		"absent":      "",
		"unparseable": "generation=banana",
		"negative":    "generation=-1",
	} {
		rr := httptest.NewRecorder()
		AgentSync(db)(rr, agentReq(t, http.MethodGet, "/agent/sync?"+q, nil, agent, nil))
		if rr.Code != http.StatusOK {
			t.Errorf("%s = %d, want 200 — an agent that cannot say where it is must be re-delivered to", name, rr.Code)
		}
	}
}

// TestOnlyATerminalResultReleasesTheSlot.
//
// What breaks in production if this fails: either a deadlock or a double run,
// and the two failure directions have opposite causes.
//
// If assignment stops restating an in-flight task, the agent that lost a
// response — or restarted between poll and start — is handed a *second* task
// while the first container is still running, and two make targets execute
// against one cluster concurrently. Nothing in the control plane will report
// that; both will look like ordinary runs.
//
// If something other than a terminal report releases the slot (a lease
// expiring, a start being treated as an acknowledgement), the same thing
// happens on a slow container. And if a terminal report fails to release it,
// the agent is stuck on a finished task forever and its cluster silently stops
// accepting work.
func TestOnlyATerminalResultReleasesTheSlot(t *testing.T) {
	db := pgtest.DB(t)

	org, cluster, server := seedAgentWorld(t, db)
	agent := seedEnrolledAgent(t, db, org, cluster, server)

	run1 := seedQueuedRun(t, db, org, cluster, "bootstrap")
	run2 := seedQueuedRun(t, db, org, cluster, "upgrade")

	first, ok := assignOnce(t, db, agent)
	if !ok || first.State != models.AgentTaskStateAssigned {
		t.Fatalf("no initial assignment: %+v", first)
	}
	params := map[string]string{"taskID": first.TaskID}

	restated := func(stage string, wantState string) {
		t.Helper()
		for i := 0; i < 2; i++ {
			as, ok := assignOnce(t, db, agent)
			if !ok || as.TaskID != first.TaskID {
				t.Fatalf("%s: slot released early, got %+v (ok=%v)", stage, as, ok)
			}
			if as.State != wantState {
				t.Fatalf("%s: state = %q, want %q", stage, as.State, wantState)
			}
		}
	}

	restated("assigned", models.AgentTaskStateAssigned)

	rr := httptest.NewRecorder()
	AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
		agentTaskStartRequest{ContainerID: "c0ffee"}, agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("start = %d %s", rr.Code, rr.Body.String())
	}
	// Start records a container id and nothing else. It is not an
	// acknowledgement and it releases nothing.
	restated("started", models.AgentTaskStateStarted)

	rr = httptest.NewRecorder()
	AgentTaskLogs(db)(rr, agentReq(t, http.MethodPost, "/logs", agentTaskLogsRequest{
		Chunks: []agentLogChunk{{Seq: 1, Stream: "stdout", Chunk: "working"}},
	}, agent, params))
	if rr.Code != http.StatusOK {
		t.Fatalf("logs = %d %s", rr.Code, rr.Body.String())
	}
	restated("after logs", models.AgentTaskStateStarted)

	rr = httptest.NewRecorder()
	AgentReconcileReport(db)(rr, agentReq(t, http.MethodPost, "/reconcile-report",
		agentReconcileReportRequest{CurrentGeneration: 0, AppliedGeneration: 0, Healthy: true}, agent, nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("reconcile = %d %s", rr.Code, rr.Body.String())
	}
	restated("after reconcile report", models.AgentTaskStateStarted)

	// Only now, and a failure counts: the slot is released by terminality, not
	// by success.
	rr = httptest.NewRecorder()
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{State: models.AgentTaskStateFailed, Error: "make exploded"}, agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("result = %d %s", rr.Code, rr.Body.String())
	}

	second, ok := assignOnce(t, db, agent)
	if !ok {
		t.Fatal("slot never released — this agent's cluster now accepts no work at all")
	}
	if second.TaskID == first.TaskID || second.Command != "upgrade" {
		t.Fatalf("second assignment = %+v, want a fresh task for run2", second)
	}
	var task2 models.AgentTask
	if err := db.First(&task2, "id = ?", second.TaskID).Error; err != nil {
		t.Fatal(err)
	}
	if task2.RunID != run2.ID {
		t.Fatalf("second task points at %v, want run2 %v", task2.RunID, run2.ID)
	}

	var settled models.ClusterRun
	if err := db.First(&settled, "id = ?", run1.ID).Error; err != nil {
		t.Fatal(err)
	}
	if settled.Status != models.ClusterRunStatusFailed || settled.Error != "make exploded" {
		t.Fatalf("run1 = %q/%q, want failed with the agent's error", settled.Status, settled.Error)
	}
}

// TestDeadLetteredTaskFailsItsRunAndIsNeverRedriven.
//
// What breaks in production if this fails: the dead letter is the mandatory
// third exit, and it has two halves that must both hold.
//
// If it does not settle the parent run, the run sits in `running` forever. The
// log reader defines terminality negatively, so it reports done:false
// indefinitely and the cluster page shows a spinner that never resolves — for a
// container nobody can say anything about. That is precisely the limbo the
// third outcome exists to prevent, and it looks like a hang rather than a bug.
//
// If anything redrives it, a make target that may have half-applied is re-run
// unattended. Recovery is a human posting a new run, which produces a new run
// row and therefore a new task; nothing in this path may do that on its own.
func TestDeadLetteredTaskFailsItsRunAndIsNeverRedriven(t *testing.T) {
	db := pgtest.DB(t)

	org, cluster, server := seedAgentWorld(t, db)
	agent := seedEnrolledAgent(t, db, org, cluster, server)
	run := seedQueuedRun(t, db, org, cluster, "bootstrap")

	as, ok := assignOnce(t, db, agent)
	if !ok {
		t.Fatal("no assignment")
	}
	params := map[string]string{"taskID": as.TaskID}

	rr := httptest.NewRecorder()
	AgentTaskStart(db)(rr, agentReq(t, http.MethodPost, "/start",
		agentTaskStartRequest{ContainerID: "abandoned"}, agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("start = %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	AgentTaskResult(db)(rr, agentReq(t, http.MethodPost, "/result",
		agentTaskResultRequest{
			State: models.AgentTaskStateDeadLettered,
			Error: "container vanished; cannot determine outcome",
		}, agent, params))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("dead letter = %d %s", rr.Code, rr.Body.String())
	}

	var task models.AgentTask
	if err := db.First(&task, "id = ?", as.TaskID).Error; err != nil {
		t.Fatal(err)
	}
	if task.State != models.AgentTaskStateDeadLettered {
		t.Fatalf("task state = %q", task.State)
	}
	// The reason is what tells an operator who has to act, and it is the only
	// place the distinction survives — the run below is a plain failure.
	if task.DeadLetterReason != models.AgentTaskDeadLetterAgentReported {
		t.Fatalf("dead_letter_reason = %q, want %q", task.DeadLetterReason, models.AgentTaskDeadLetterAgentReported)
	}
	if task.EndedAt == nil {
		t.Fatal("ended_at not set on a terminal task")
	}

	var settled models.ClusterRun
	if err := db.First(&settled, "id = ?", run.ID).Error; err != nil {
		t.Fatal(err)
	}
	if settled.Status != models.ClusterRunStatusFailed {
		t.Fatalf("run status = %q, want failed — the log reader will tail this forever", settled.Status)
	}
	if settled.FinishedAt.IsZero() {
		t.Fatal("finished_at not set on a settled run")
	}
	if !strings.Contains(settled.Error, models.AgentTaskDeadLetterAgentReported) ||
		!strings.Contains(settled.Error, "cannot determine outcome") {
		t.Fatalf("run error = %q, want the reason and the agent's message", settled.Error)
	}

	// No redrive: the slot is free, and polling it repeatedly must not produce
	// another attempt at the same run.
	for i := 0; i < 3; i++ {
		if again, ok := assignOnce(t, db, agent); ok {
			t.Fatalf("poll %d redrove a dead letter: %+v", i, again)
		}
	}
	var n int64
	db.Model(&models.AgentTask{}).Where("run_id = ?", run.ID).Count(&n)
	if n != 1 {
		t.Fatalf("%d tasks for a dead-lettered run, want 1", n)
	}
}

// TestAgentLogRedeliveryDoesNotDuplicateRows.
//
// What breaks in production if this fails: the log tail becomes untrustworthy
// and nobody can tell. The agent is an at-least-once client — it redelivers a
// batch whenever a response is lost — and job_logs has no unique key to catch
// it, by design: the dedupe is the per-task watermark. Lose it and a flaky
// network turns a bootstrap log into stuttering repeated output, which reads as
// the make target genuinely looping. Retention then reaps a multiple of what it
// was sized for.
//
// The partial case is the one worth having: the server writing rows and dying
// before it answers is the normal failure, and the agent's next batch overlaps
// what already landed.
func TestAgentLogRedeliveryDoesNotDuplicateRows(t *testing.T) {
	db := pgtest.DB(t)

	org, cluster, server := seedAgentWorld(t, db)
	agent := seedEnrolledAgent(t, db, org, cluster, server)
	run := seedQueuedRun(t, db, org, cluster, "bootstrap")

	as, ok := assignOnce(t, db, agent)
	if !ok {
		t.Fatal("no assignment")
	}
	params := map[string]string{"taskID": as.TaskID}

	post := func(seqs ...int64) int64 {
		t.Helper()
		chunks := make([]agentLogChunk, 0, len(seqs))
		for _, s := range seqs {
			chunks = append(chunks, agentLogChunk{
				Seq: s, Stream: "stdout", Chunk: "line-" + strconv.FormatInt(s, 10),
			})
		}
		rr := httptest.NewRecorder()
		AgentTaskLogs(db)(rr, agentReq(t, http.MethodPost, "/logs",
			agentTaskLogsRequest{Chunks: chunks}, agent, params))
		if rr.Code != http.StatusOK {
			t.Fatalf("logs %v = %d %s", seqs, rr.Code, rr.Body.String())
		}
		var res agentTaskLogsResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			t.Fatal(err)
		}
		return res.AcceptedThroughSeq
	}

	countRows := func() int64 {
		t.Helper()
		var n int64
		db.Model(&models.JobLog{}).
			Where("subject_type = ? AND subject_id = ?", models.JobLogSubjectClusterRun, run.ID).
			Count(&n)
		return n
	}

	if got := post(1, 2, 3); got != 3 {
		t.Fatalf("accepted = %d, want 3", got)
	}
	if got := countRows(); got != 3 {
		t.Fatalf("rows = %d, want 3", got)
	}

	// Whole-batch redelivery.
	if got := post(1, 2, 3); got != 3 {
		t.Fatalf("redelivery accepted = %d, want the unchanged watermark 3", got)
	}
	if got := countRows(); got != 3 {
		t.Fatalf("rows after full redelivery = %d, want 3", got)
	}

	// A batch entirely below the watermark writes nothing and does not move it
	// backwards. A watermark that regressed here would re-admit everything
	// after it on the next batch.
	if got := post(2); got != 3 {
		t.Fatalf("stale batch accepted = %d, want 3", got)
	}
	if got := countRows(); got != 3 {
		t.Fatalf("rows after stale batch = %d, want 3", got)
	}

	// The real recovery case: overlapping, partly new.
	if got := post(2, 3, 4, 5); got != 5 {
		t.Fatalf("overlapping accepted = %d, want 5", got)
	}
	if got := countRows(); got != 5 {
		t.Fatalf("rows after overlap = %d, want 5 (2 and 3 duplicated)", got)
	}

	var task models.AgentTask
	if err := db.First(&task, "id = ?", as.TaskID).Error; err != nil {
		t.Fatal(err)
	}
	if task.LogSeq != 5 {
		t.Fatalf("log_seq = %d, want 5", task.LogSeq)
	}

	// The existing reader has to see exactly the same thing. Agent chunks land
	// under subject_type='cluster_run' with the run's id precisely so no reader
	// can tell an agent wrote them rather than a River worker.
	items, _, err := readJobLogs(db, org.ID, models.JobLogSubjectClusterRun, run.ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 5 {
		t.Fatalf("readJobLogs = %d rows, want 5", len(items))
	}
	for i, want := range []string{"line-1", "line-2", "line-3", "line-4", "line-5"} {
		if items[i].Chunk != want {
			t.Fatalf("chunk %d = %q, want %q — the tail is out of order or duplicated", i, items[i].Chunk, want)
		}
	}
}

// TestAgentReconcileStatusCannotOverwriteAnotherAgents.
//
// What breaks in production if this fails: convergence reporting lies. The
// upsert key is (agent_id, resource_type, resource_id), and agent_id is the
// only part of it that is not attacker-supplied — the resource type and id come
// straight out of the request body. Drop agent_id from the conflict target and
// any agent can overwrite any other's status for a resource of the same name,
// so a broken cluster reports itself applied because a healthy one said so.
func TestAgentReconcileStatusCannotOverwriteAnotherAgents(t *testing.T) {
	db := pgtest.DB(t)

	orgA, clusterA, serverA := seedAgentWorld(t, db)
	agentA := seedEnrolledAgent(t, db, orgA, clusterA, serverA)
	orgB, clusterB, serverB := seedAgentWorld(t, db)
	agentB := seedEnrolledAgent(t, db, orgB, clusterB, serverB)

	report := func(a *models.Agent, state, lastErr string) {
		t.Helper()
		rr := httptest.NewRecorder()
		AgentReconcileReport(db)(rr, agentReq(t, http.MethodPost, "/reconcile-report",
			agentReconcileReportRequest{
				CurrentGeneration: 7, AppliedGeneration: 7, Healthy: state == models.AgentReconcileStateApplied,
				Resources: []agentResourceStatus{{
					// Same resource identity on both sides on purpose: these
					// strings are the agent's to choose.
					ResourceType: "kubeconfig", ResourceID: "default",
					DesiredGeneration: 7, State: state, LastError: lastErr,
				}},
			}, a, nil))
		if rr.Code != http.StatusNoContent {
			t.Fatalf("reconcile = %d %s", rr.Code, rr.Body.String())
		}
	}

	report(agentB, models.AgentReconcileStateFailed, "cluster is on fire")
	report(agentA, models.AgentReconcileStateApplied, "")

	var b models.AgentReconcileStatus
	if err := db.Where("agent_id = ?", agentB.ID).First(&b).Error; err != nil {
		t.Fatal(err)
	}
	if b.State != models.AgentReconcileStateFailed || b.LastError != "cluster is on fire" {
		t.Fatalf("agent B's status was overwritten by agent A: %+v", b)
	}
	if b.ClusterID != clusterB.ID || b.OrganizationID != orgB.ID {
		t.Fatalf("status carries the wrong scope: %+v", b)
	}

	var a models.AgentReconcileStatus
	if err := db.Where("agent_id = ?", agentA.ID).First(&a).Error; err != nil {
		t.Fatal(err)
	}
	if a.ClusterID != clusterA.ID || a.OrganizationID != orgA.ID {
		t.Fatalf("status carries the wrong scope: %+v", a)
	}

	// The scope is taken from the credential, never from the body — there is no
	// field in the request that could say otherwise, and that is the point.
	var reloaded models.Agent
	if err := db.First(&reloaded, "id = ?", agentB.ID).Error; err != nil {
		t.Fatal(err)
	}
	if reloaded.Healthy {
		t.Fatal("agent B was marked healthy by agent A's report")
	}
}
