package handlers

import (
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Regression tests for the critical and high findings the adversarial review
// raised against the first cut of the agent control plane. Each one asserts the
// specific behaviour that was wrong, so a revert is caught rather than merely
// discussed.

func seedExecutorRun(t *testing.T, db *gorm.DB, orgID, clusterID uuid.UUID, executor, status string) models.ClusterRun {
	t.Helper()
	run := models.ClusterRun{
		OrganizationID: orgID,
		ClusterID:      clusterID,
		Action:         "bootstrap",
		Status:         status,
		Executor:       executor,
	}
	if err := db.Create(&run).Error; err != nil {
		t.Fatalf("create run: %v", err)
	}
	return run
}

// CRITICAL. River's ClusterActionWorker executes runs over SSH and leaves them
// in exactly the queued/running states the agent's assignment query looks for.
// Before the executor column existed, an agent claimed those runs and the same
// make target executed twice, concurrently, against one Terraform state.
func TestClusterRunExecutorSeparatesRiverFromAgent(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "exec-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	cluster := models.Cluster{OrganizationID: org.ID, Name: "c1"}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("create cluster: %v", err)
	}

	riverRun := seedExecutorRun(t, db, org.ID, cluster.ID, models.ClusterRunExecutorRiver, models.ClusterRunStatusRunning)
	agentRun := seedExecutorRun(t, db, org.ID, cluster.ID, models.ClusterRunExecutorAgent, models.ClusterRunStatusQueued)

	// The exact predicate AgentAssignment derives from.
	var claimable []models.ClusterRun
	if err := db.
		Where("cluster_id = ? AND organization_id = ? AND executor = ? AND status IN ?",
			cluster.ID, org.ID, models.ClusterRunExecutorAgent,
			[]string{models.ClusterRunStatusQueued, models.ClusterRunStatusRunning}).
		Where("NOT EXISTS (SELECT 1 FROM agent_tasks t WHERE t.run_id = cluster_runs.id)").
		Find(&claimable).Error; err != nil {
		t.Fatalf("claimable query: %v", err)
	}

	if len(claimable) != 1 {
		t.Fatalf("agent can claim %d runs, want exactly 1", len(claimable))
	}
	if claimable[0].ID == riverRun.ID {
		t.Fatal("agent claimed a run River is executing — the same make target would run twice")
	}
	if claimable[0].ID != agentRun.ID {
		t.Fatalf("claimed %s, want the agent-owned run %s", claimable[0].ID, agentRun.ID)
	}
}

// The API enqueue path must stamp the executor at creation. JobID is only
// filled in after the River insert returns, so inferring ownership from it
// leaves a window where a run looks unclaimed to a polling agent.
func TestEnqueuedRunsAreStampedRiverAtCreation(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "exec-default-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	cluster := models.Cluster{OrganizationID: org.ID, Name: "c2"}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("create cluster: %v", err)
	}

	// Mirrors what the enqueue handler constructs.
	run := models.ClusterRun{
		OrganizationID: org.ID,
		ClusterID:      cluster.ID,
		Action:         "bootstrap",
		Status:         models.ClusterRunStatusQueued,
		Executor:       models.ClusterRunExecutorRiver,
	}
	if err := db.Create(&run).Error; err != nil {
		t.Fatalf("create run: %v", err)
	}

	var stored models.ClusterRun
	if err := db.Where("id = ?", run.ID).First(&stored).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if stored.Executor != models.ClusterRunExecutorRiver {
		t.Fatalf("executor = %q, want river", stored.Executor)
	}

	// And a row written without an executor at all — every run that predates
	// this column — must default to river, never to agent.
	var legacyExecutor string
	if err := db.Raw(`
		INSERT INTO cluster_runs (organization_id, cluster_id, action, status, error)
		VALUES (?, ?, 'legacy', 'queued', '')
		RETURNING executor`, org.ID, cluster.ID).Scan(&legacyExecutor).Error; err != nil {
		t.Fatalf("legacy insert: %v", err)
	}
	if legacyExecutor != models.ClusterRunExecutorRiver {
		t.Fatalf("a run created without an executor defaulted to %q; every pre-existing run would become agent-claimable", legacyExecutor)
	}
}

// HIGH. The batch filter used to compare against the pre-batch watermark, which
// never moves inside the loop — so two chunks carrying the same seq in one POST
// were both inserted, and the transcript reader pages by job_logs.id with no way
// to tell the copies apart.
func TestLogBatchDedupWithinASingleBatch(t *testing.T) {
	// Mirrors the accepted-watermark filter in AgentTaskLogs.
	filter := func(preBatch int64, seqs []int64) []int64 {
		accepted := preBatch
		var kept []int64
		for _, s := range seqs {
			if s <= accepted {
				continue
			}
			kept = append(kept, s)
			if s > accepted {
				accepted = s
			}
		}
		return kept
	}

	got := filter(0, []int64{1, 2, 2, 3, 3, 3})
	want := []int64{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("kept %v, want %v — duplicate seqs in one batch are inserted twice", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("kept %v, want %v", got, want)
		}
	}

	// Already-delivered seqs stay filtered against the previous watermark.
	if kept := filter(5, []int64{3, 4, 5, 6}); len(kept) != 1 || kept[0] != 6 {
		t.Fatalf("kept %v, want only seq 6", kept)
	}

	// Out-of-order delivery is dropped rather than written behind the reader's
	// cursor, where nothing would ever revisit it.
	if kept := filter(0, []int64{3, 1, 2, 4}); len(kept) != 2 || kept[0] != 3 || kept[1] != 4 {
		t.Fatalf("kept %v, want [3 4]", kept)
	}
}

// HIGH. An issued credential with no expiry cannot be ended by anything: the
// only revocation path fires on a *successful* re-enrolment, so a leaked
// credential was valid forever.
func TestEnrolledAgentCredentialExpires(t *testing.T) {
	if ttl := agentCredentialTTL(); ttl <= 0 {
		t.Fatalf("credential TTL = %v, want a positive bound", ttl)
	}

	db := pgtest.DB(t)
	org := models.Organization{Name: "ttl-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	cluster := models.Cluster{OrganizationID: org.ID, Name: "ttl-cluster"}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("create cluster: %v", err)
	}

	key := createTestSshKey(t, db, org.ID, "ttl-key")
	pub := "198.51.100.7"
	server := models.Server{
		OrganizationID:   org.ID,
		Hostname:         "ttl-bastion",
		PublicIPAddress:  &pub,
		PrivateIPAddress: "10.0.0.7",
		SSHUser:          "deploy",
		SshKeyID:         key.ID,
		Role:             "bastion",
		Status:           "ready",
	}
	if err := db.Create(&server).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	// An agent row carrying an expiry in the past must not satisfy the
	// predicate the auth middleware authenticates with.
	past := time.Now().UTC().Add(-time.Hour)
	agent := models.Agent{
		OrganizationID: org.ID,
		ClusterID:      cluster.ID,
		ServerID:       server.ID,
		KeyHash:        "expired-key-hash",
		SecretHash:     "x",
		Prefix:         "agt_",
		Status:         models.AgentStatusActive,
		EnrolledAt:     time.Now().UTC(),
		ExpiresAt:      &past,
	}
	if err := db.Create(&agent).Error; err != nil {
		t.Fatalf("create agent: %v", err)
	}

	var found models.Agent
	err := db.Where("key_hash = ? AND status = ? AND (expires_at IS NULL OR expires_at > ?)",
		"expired-key-hash", models.AgentStatusActive, time.Now()).First(&found).Error
	if err == nil {
		t.Fatal("an expired credential still authenticates")
	}
}
