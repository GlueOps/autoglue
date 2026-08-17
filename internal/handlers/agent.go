package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// The agent plane is deliberately absent from docs/openapi.json. No handler in
// this file carries a swag annotation block, and a route annotation is the only
// thing swag looks for, so omission is the whole mechanism — nothing has to be
// excluded anywhere else. (The token is spelled out nowhere in this file on
// purpose: swag scans comments, not just the ones attached to a handler.)
// The corollary is that the DTOs below are unexported plain structs rather than
// members of internal/handlers/dto: they will never reach the embedded spec or
// the generated SDKs, so nothing may be built on the assumption that they do.

// Poll intervals handed back to the agent. Advisory rather than enforced: the
// control plane would rather state a cadence it can absorb than have every
// bastion pick its own and discover the aggregate at rollout.
const (
	agentSyncPollSec       = 30
	agentAssignmentPollSec = 10
)

// defaultAgentCredentialTTLDays bounds how long a leaked agent credential stays
// useful. Long, because re-enrolment is currently a manual push over SSH and an
// expiry shorter than the interval at which anyone touches a bastion would take
// the fleet down; shorten it via agent.credential_ttl_days once re-enrolment is
// automated.
const defaultAgentCredentialTTLDays = 90

func agentCredentialTTL() time.Duration {
	if d := viper.GetInt("agent.credential_ttl_days"); d > 0 {
		return time.Duration(d) * 24 * time.Hour
	}
	return defaultAgentCredentialTTLDays * 24 * time.Hour
}

func ptrTime(t time.Time) *time.Time { return &t }

// clusterRunStatusSucceeded is the literal the rest of the system already
// agrees on. models.ClusterRunStatusSuccess is "success" and is dead code —
// bg/cluster_action.go writes "succeeded" and the UI switches on "succeeded" —
// so using the constant here would introduce a third spelling of the same
// state. Reconciling the constant is a models change and out of scope for this
// file; writing what the existing readers understand is not optional.
const clusterRunStatusSucceeded = "succeeded"

// agentLogStreamStderr is docker's second stream. models/job_log.go declares
// only stdout and system, because runMakeOnBastion merges the two before
// writing; an agent following a container does not, so the distinction survives
// to here and is worth keeping. Belongs in that const block once models is open
// for edit again.
const agentLogStreamStderr = "stderr"

// Bounds on a single log POST. The global 10 MiB body cap would eventually stop
// an abusive batch, but as a truncated-JSON parse error rather than something an
// agent author can act on, so the limits are stated here instead.
const (
	agentMaxLogChunkBytes = 512 << 10
	agentMaxLogChunks     = 512
)

// ---------------------------------------------------------------------------
// wire types
// ---------------------------------------------------------------------------

type agentEnrollRequest struct {
	Ticket  string `json:"ticket"`
	Version string `json:"version"`
}

type agentEnrollResponse struct {
	AgentID        uuid.UUID `json:"agent_id"`
	Key            string    `json:"key"`
	Secret         string    `json:"secret"`
	ClusterID      uuid.UUID `json:"cluster_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	PollAfterSec   int       `json:"poll_after_sec"`
}

// agentResource mirrors the agent's api.Resource exactly, including the fact
// that DependsOn and Spec are JSON *strings* rather than nested objects: the
// agent stores them verbatim in a SQLite TEXT column and only parses them when
// it applies the resource. Re-encoding them as objects here would force it to
// re-serialise, and the spec hash is computed control-plane side precisely so
// both ends agree on the bytes rather than on how a marshaller happened to
// order them.
type agentResource struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	ResourceName string `json:"resource_name"`
	Phase        int    `json:"phase"`
	DependsOn    string `json:"depends_on"`
	Spec         string `json:"spec"`
	SpecHash     string `json:"spec_hash"`
	Required     bool   `json:"required"`
}

type agentSnapshot struct {
	Generation   int64           `json:"generation"`
	FullSnapshot bool            `json:"full_snapshot"`
	PollAfterSec int             `json:"poll_after_sec"`
	Resources    []agentResource `json:"resources"`
}

type agentAssignmentResponse struct {
	TaskID       string `json:"task_id"`
	Command      string `json:"command"`
	Args         string `json:"args"`
	State        string `json:"state"`
	PollAfterSec int    `json:"poll_after_sec"`
}

// agentTaskArgs is the agent's supervisor.TaskArgs. It describes exactly one
// container because a task is exactly one container.
type agentTaskArgs struct {
	ClusterID string   `json:"cluster_id"`
	RunID     string   `json:"run_id"`
	Image     string   `json:"image"`
	Tag       string   `json:"tag"`
	Target    string   `json:"target"`
	WorkDir   string   `json:"work_dir"`
	Mounts    []string `json:"mounts"` // "source:target"
}

type agentTaskStartRequest struct {
	ContainerID string `json:"container_id"`
}

type agentLogChunk struct {
	Seq    int64  `json:"seq"`
	Stream string `json:"stream"`
	Chunk  string `json:"chunk"`
}

type agentTaskLogsRequest struct {
	Chunks []agentLogChunk `json:"chunks"`
}

type agentTaskLogsResponse struct {
	TaskID             uuid.UUID `json:"task_id"`
	AcceptedThroughSeq int64     `json:"accepted_through_seq"`
}

type agentTaskResultRequest struct {
	State    string `json:"state"`
	ExitCode *int   `json:"exit_code,omitempty"`
	Error    string `json:"error,omitempty"`
}

type agentResourceStatus struct {
	ResourceType      string `json:"resource_type"`
	ResourceID        string `json:"resource_id"`
	DesiredGeneration int64  `json:"desired_generation"`
	State             string `json:"state"`
	LastError         string `json:"last_error,omitempty"`
}

type agentReconcileReportRequest struct {
	CurrentGeneration int64                 `json:"current_generation"`
	AppliedGeneration int64                 `json:"applied_generation"`
	Healthy           bool                  `json:"healthy"`
	Resources         []agentResourceStatus `json:"resources,omitempty"`
}

// ---------------------------------------------------------------------------
// enrolment
// ---------------------------------------------------------------------------

// EnrollAgent redeems a single-use enrolment ticket and returns the credential
// tuple the agent will authenticate with from then on.
//
// This is the one agent endpoint that cannot be behind AgentAuth: it is the
// thing that issues the credential AgentAuth checks. The ticket in the body is
// the credential, and it is trustworthy because of how it arrived — pushed to
// the bastion over SSH alongside the rest of the cluster assets — so redeeming
// it demonstrates the caller reached that host.
//
// Every rejection is the same 401 with the same body. A ticket that is unknown,
// a ticket that expired and a ticket already redeemed are three different facts
// about a secret the caller is holding, and distinguishing them turns this
// endpoint into an oracle for whether a leaked ticket is still worth trying.
func EnrollAgent(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req agentEnrollRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid body")
			return
		}
		ticketPlain := strings.TrimSpace(req.Ticket)
		if ticketPlain == "" {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid enrollment ticket")
			return
		}

		var (
			agent models.Agent
			cred  auth.AgentCredential
		)
		txErr := db.Transaction(func(tx *gorm.DB) error {
			now := time.Now().UTC()

			// FOR UPDATE plus the redeemed_at/expires_at predicates makes the
			// check and the claim one statement: two bastions racing the same
			// ticket serialise here, and the loser finds redeemed_at set.
			var ticket models.AgentEnrollmentTicket
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("ticket_hash = ? AND redeemed_at IS NULL AND expires_at > ?",
					auth.SHA256Hex(ticketPlain), now).
				First(&ticket).Error; err != nil {
				return err
			}

			// Minting happens here, after the ticket has been proven, and never
			// before. MintAgentCredential runs argon2id at 64 MiB, so hoisting
			// it above the lookup hands any unauthenticated caller a
			// memory-hard grinder they can drive with a garbage ticket.
			var err error
			if cred, err = auth.MintAgentCredential(); err != nil {
				return err
			}

			// The prior agent has to be revoked before the new row is inserted,
			// not after: uniq_agents_live_cluster is what guarantees one live
			// agent per cluster, so an insert ahead of the revoke is the one
			// ordering the database will refuse.
			if err := revokeClusterAgents(tx, ticket.ClusterID, now); err != nil {
				return err
			}

			agent = models.Agent{
				ID:             cred.ID,
				OrganizationID: ticket.OrganizationID,
				ClusterID:      ticket.ClusterID,
				ServerID:       ticket.ServerID,
				KeyHash:        cred.KeyHash,
				SecretHash:     cred.SecretHash,
				Prefix:         cred.Prefix,
				Status:         models.AgentStatusActive,
				Version:        strings.TrimSpace(req.Version),
				EnrolledAt:     now,
				// Bounded on purpose. Re-enrolment revokes the previous agent,
				// but that only fires when someone enrols again — so without an
				// expiry a leaked credential stays valid forever and nothing in
				// the system can end it.
				//
				// The counter-argument is real: there is no renewal endpoint, so
				// an expiry is a date at which a healthy agent stops working.
				// That is survivable only because re-enrolment is the recovery
				// path — mint a ticket, push it over SSH, the agent enrols again
				// — which is the same path used when a bastion is rebuilt. The
				// TTL is therefore long relative to how often that path runs, and
				// configurable so it can be shortened once the sweeper that
				// automates re-enrolment lands.
				ExpiresAt: ptrTime(now.Add(agentCredentialTTL())),
			}
			if err := tx.Create(&agent).Error; err != nil {
				return err
			}

			return tx.Model(&models.AgentEnrollmentTicket{}).
				Where("id = ?", ticket.ID).
				Updates(map[string]any{
					"redeemed_at":          now,
					"redeemed_by_agent_id": agent.ID,
				}).Error
		})
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid enrollment ticket")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// The only time the plaintext exists outside the bastion. Nothing
		// persisted it; a lost response means a new ticket, not a lookup.
		utils.WriteJSON(w, http.StatusCreated, agentEnrollResponse{
			AgentID:        agent.ID,
			Key:            cred.Key,
			Secret:         cred.Secret,
			ClusterID:      agent.ClusterID,
			OrganizationID: agent.OrganizationID,
			PollAfterSec:   agentSyncPollSec,
		})
	}
}

// revokeClusterAgents retires every live agent on a cluster and settles whatever
// they were holding.
//
// Dead-lettering the in-flight task is not tidying up. The partial unique index
// counts an assigned or started task against the cluster whether or not its
// agent still exists, so leaving it would block the new agent from ever being
// given work, and its ClusterRun would sit in `running` forever with the log
// reader reporting done:false — the exact limbo the third outcome exists to
// prevent. Reason `agent_reenrolled` rather than a plain failure because the
// distinction decides who has to act: nobody can say whether that container
// finished, so a human decides whether to re-run it.
func revokeClusterAgents(tx *gorm.DB, clusterID uuid.UUID, now time.Time) error {
	var prior []models.Agent
	if err := tx.Where("cluster_id = ? AND status = ?", clusterID, models.AgentStatusActive).
		Find(&prior).Error; err != nil {
		return err
	}

	for _, p := range prior {
		var held []models.AgentTask
		if err := tx.Where("agent_id = ? AND state IN ?", p.ID,
			[]string{models.AgentTaskStateAssigned, models.AgentTaskStateStarted}).
			Find(&held).Error; err != nil {
			return err
		}

		for _, t := range held {
			msg := "agent re-enrolled while this task was in flight"
			if err := tx.Model(&models.AgentTask{}).
				Where("id = ?", t.ID).
				Updates(map[string]any{
					"state":              models.AgentTaskStateDeadLettered,
					"dead_letter_reason": models.AgentTaskDeadLetterReenrolled,
					"error":              msg,
					"ended_at":           now,
				}).Error; err != nil {
				return err
			}
			if err := failRunFromDeadLetter(tx, t.RunID,
				models.AgentTaskDeadLetterReenrolled, msg, now); err != nil {
				return err
			}
		}

		if err := tx.Model(&models.Agent{}).
			Where("id = ?", p.ID).
			Updates(map[string]any{
				"status":     models.AgentStatusRevoked,
				"revoked_at": now,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

// failRunFromDeadLetter settles the ClusterRun behind a dead-lettered task.
//
// The run is failed rather than given a run-level dead_lettered status because
// ClusterRun.Status has live consumers that enumerate its values — the cluster
// page and the log reader — and neither handles a sixth one. The distinction is
// not lost: it survives on the task as state plus dead_letter_reason, which is
// where an operator or a future admin endpoint reads it from.
//
// The status predicate is what makes a redelivered result harmless: a run that
// has already settled cannot be resurrected by a late report.
func failRunFromDeadLetter(tx *gorm.DB, runID uuid.UUID, reason, taskErr string, now time.Time) error {
	msg := "dead-lettered (" + reason + ")"
	if taskErr != "" {
		msg += ": " + taskErr
	}
	return tx.Model(&models.ClusterRun{}).
		Where("id = ? AND status IN ?", runID,
			[]string{models.ClusterRunStatusQueued, models.ClusterRunStatusRunning}).
		Updates(map[string]any{
			"status":      models.ClusterRunStatusFailed,
			"error":       msg,
			"finished_at": now,
		}).Error
}

// ---------------------------------------------------------------------------
// config plane
// ---------------------------------------------------------------------------

// AgentSync delivers the cluster's whole desired-state snapshot, or 204 when the
// caller already holds the current generation.
//
// A GET even though the agent will rewrite its disk from the response: the
// mutation lands on the bastion, not on control-plane state, which is the same
// reason fetching a firmware image is not a POST. The generation travels as a
// query parameter rather than a body because a GET body has no defined
// semantics and is dropped by some proxies — and a dropped generation reads as
// zero, silently turning every poll into a full re-delivery.
//
// Always a whole snapshot, never a diff. An agent that was switched off for a
// week has no history to apply a diff against, and a snapshot is correct from
// any starting state including a wiped disk.
func AgentSync(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, ok := httpmiddleware.AgentFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "agent required")
			return
		}

		// An unparseable or absent generation is treated as zero rather than
		// rejected: "I have nothing" is the honest reading of a claim the agent
		// could not state, and re-delivering a snapshot is always safe.
		claimed, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("generation")), 10, 64)
		if claimed < 0 {
			claimed = 0
		}

		var desired models.ClusterDesiredState
		if err := db.Where("cluster_id = ?", agent.ClusterID).First(&desired).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Nothing has ever been published for this cluster, so the agent
				// is trivially current. 204 rather than an empty snapshot: an
				// empty resource list at generation zero would read as "delete
				// everything you have".
				w.WriteHeader(http.StatusNoContent)
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if desired.Generation <= claimed {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Selecting one generation by equality is what makes this safe to serve
		// without a lock: rows for a generation are written once and never
		// mutated, and the pointer is bumped in the same transaction that wrote
		// them, so a reader cannot see a pointer aimed at a half-written set.
		var rows []models.DesiredResource
		if err := db.
			Where("cluster_id = ? AND generation = ?", agent.ClusterID, desired.Generation).
			Order("phase ASC, resource_id ASC").
			Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := agentSnapshot{
			Generation:   desired.Generation,
			FullSnapshot: true,
			PollAfterSec: agentSyncPollSec,
			Resources:    make([]agentResource, 0, len(rows)),
		}
		for _, row := range rows {
			out.Resources = append(out.Resources, agentResource{
				ResourceType: row.ResourceType,
				ResourceID:   row.ResourceID,
				ResourceName: row.ResourceName,
				Phase:        row.Phase,
				DependsOn:    jsonOrDefault(row.DependsOn, "[]"),
				Spec:         jsonOrDefault(row.Spec, "{}"),
				SpecHash:     row.SpecHash,
				Required:     row.Required,
			})
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// jsonOrDefault renders a jsonb column as the string the agent expects, falling
// back when the column is NULL. The agent parses these lazily, so handing it a
// bare "" would turn a missing dependency list into a parse error at apply time
// rather than an empty one at read time.
func jsonOrDefault(v datatypes.JSON, def string) string {
	s := strings.TrimSpace(string(v))
	if s == "" || s == "null" {
		return def
	}
	return s
}

// AgentReconcileReport records what the agent says about its own convergence.
//
// This is also the config plane's heartbeat, which is why it runs on that
// loop's timer rather than between tasks: an agent that is happily draining
// tasks while its reconcile loop is wedged should look stale here, because it
// is.
//
// Nothing in the control plane ever writes these rows on the agent's behalf. A
// status that says applied means an agent claimed it, and the gap between
// Agent.AppliedGeneration and ClusterDesiredState.Generation is the only
// convergence signal worth trusting.
func AgentReconcileReport(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, ok := httpmiddleware.AgentFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "agent required")
			return
		}

		var req agentReconcileReportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid body")
			return
		}

		now := time.Now().UTC()

		txErr := db.Transaction(func(tx *gorm.DB) error {
			if len(req.Resources) > 0 {
				rows := make([]models.AgentReconcileStatus, 0, len(req.Resources))
				for _, rs := range req.Resources {
					if rs.ResourceType == "" || rs.ResourceID == "" {
						// A status with no subject cannot be upserted onto the
						// unique key, and guessing which resource was meant is
						// worse than dropping it.
						continue
					}
					rows = append(rows, models.AgentReconcileStatus{
						OrganizationID:    agent.OrganizationID,
						ClusterID:         agent.ClusterID,
						AgentID:           agent.ID,
						ResourceType:      rs.ResourceType,
						ResourceID:        rs.ResourceID,
						DesiredGeneration: rs.DesiredGeneration,
						State:             normalizeReconcileState(rs.State),
						LastError:         rs.LastError,
						ReportedAt:        now,
					})
				}
				if len(rows) > 0 {
					// Current status only, never history: there is one truth per
					// resource and it carries the generation it was achieved at,
					// which is what makes a stale desired_generation readable as
					// "has not caught up" without diffing anything.
					if err := tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{
							{Name: "agent_id"}, {Name: "resource_type"}, {Name: "resource_id"},
						},
						DoUpdates: clause.AssignmentColumns([]string{
							"desired_generation", "state", "last_error", "reported_at", "updated_at",
						}),
					}).CreateInBatches(&rows, 200).Error; err != nil {
						return err
					}
				}
			}

			return tx.Model(&models.Agent{}).
				Where("id = ?", agent.ID).
				Updates(map[string]any{
					"reported_generation": req.CurrentGeneration,
					// Clamped monotonic. An agent that lost its disk restarts at
					// zero and would otherwise walk the pointer backwards,
					// making a converged cluster look like it regressed.
					"applied_generation": gorm.Expr("GREATEST(applied_generation, ?)", req.AppliedGeneration),
					"healthy":            req.Healthy,
					"last_reconcile_at":  now,
					"last_seen_at":       now,
				}).Error
		})
		if txErr != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// normalizeReconcileState keeps an unknown value from becoming a state the
// control plane cannot reason about. Pending is the safe fallback because it is
// the only one that claims nothing.
func normalizeReconcileState(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case models.AgentReconcileStateApplied:
		return models.AgentReconcileStateApplied
	case models.AgentReconcileStateFailed:
		return models.AgentReconcileStateFailed
	default:
		return models.AgentReconcileStatePending
	}
}

// ---------------------------------------------------------------------------
// task plane
// ---------------------------------------------------------------------------

// AgentAssignment states the task this agent should be working on, deriving it
// from durable rows rather than popping a queue.
//
// Nothing is leased and hand-out is not recorded, because hand-out proves
// nothing: a task the agent never received and a task it received and crashed
// on are indistinguishable from here. So an unacknowledged task is simply named
// again on the next poll, a lost response costs nothing, and only a reported
// terminal outcome releases the slot.
//
// Idle is 200 with a literal null, not 404. The agent decodes straight into a
// *Assignment, and an error status would make "no work" indistinguishable from
// "the control plane is broken" in its retry logic.
func AgentAssignment(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, ok := httpmiddleware.AgentFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "agent required")
			return
		}

		// 1. Restate. At most one row, guaranteed by uniq_agent_tasks_in_flight.
		// State travels with it so the agent can tell "you never started this"
		// from "you told me you started it" — drift surfaces on the next poll
		// instead of at a lease expiry that does not exist.
		var held models.AgentTask
		err := db.Where("agent_id = ? AND state IN ?", agent.ID,
			[]string{models.AgentTaskStateAssigned, models.AgentTaskStateStarted}).
			First(&held).Error
		if err == nil {
			writeAssignment(w, held)
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// 2. Derive. The oldest unclaimed run for this cluster. Only queued and
		// running are non-terminal; job_logs.go defines terminality negatively,
		// so introducing a third non-terminal status would make the log reader
		// report done mid-run.
		//
		// The executor filter is load-bearing, not defensive. River's
		// ClusterActionWorker executes runs over SSH from the control plane and
		// leaves them in exactly the same queued/running states this query looks
		// for, so without it an agent claims a run River is part-way through and
		// the same make target executes twice against one Terraform state.
		var run models.ClusterRun
		err = db.
			Where("cluster_id = ? AND organization_id = ? AND executor = ? AND status IN ?",
				agent.ClusterID, agent.OrganizationID, models.ClusterRunExecutorAgent,
				[]string{models.ClusterRunStatusQueued, models.ClusterRunStatusRunning}).
			Where("NOT EXISTS (SELECT 1 FROM agent_tasks t WHERE t.run_id = cluster_runs.id)").
			Order("created_at ASC").
			First(&run).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteJSON(w, http.StatusOK, nil)
			return
		}
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		args, err := buildTaskArgs(db, agent, run)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// 3. Materialise. Two unique indexes make this safe under concurrent
		// polls with no locking: run_id uniqueness stops a second task being
		// minted for one run, and uniq_agent_tasks_in_flight stops a second task
		// for this agent even when two runs are queued. A losing insert is
		// therefore "no assignment this tick", not an error — the task row
		// exists because a run exists, not because a poll happened.
		task := models.AgentTask{
			OrganizationID: agent.OrganizationID,
			ClusterID:      agent.ClusterID,
			AgentID:        agent.ID,
			RunID:          run.ID,
			Kind:           models.AgentTaskKindContainer,
			Command:        run.Action,
			Args:           args,
			State:          models.AgentTaskStateAssigned,
			AssignedAt:     time.Now().UTC(),
		}
		res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&task)
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if res.RowsAffected == 0 {
			// Another poll got there first. Saying nothing is correct: the next
			// poll restates whatever actually won, and derivation is a pure
			// function so it will agree.
			utils.WriteJSON(w, http.StatusOK, nil)
			return
		}

		// Advisory for the UI and the log reader, which is why a failure here is
		// not fatal: the task is the durable record of the work, and a run left
		// in `queued` is still non-terminal so the reader keeps tailing.
		_ = db.Model(&models.ClusterRun{}).
			Where("id = ? AND status = ?", run.ID, models.ClusterRunStatusQueued).
			Update("status", models.ClusterRunStatusRunning).Error

		writeAssignment(w, task)
	}
}

func writeAssignment(w http.ResponseWriter, t models.AgentTask) {
	utils.WriteJSON(w, http.StatusOK, agentAssignmentResponse{
		TaskID:  t.ID.String(),
		Command: t.Command,
		// Args is a JSON string rather than a nested object because that is what
		// the agent stores and re-parses; see agentResource for the same rule.
		Args:         jsonOrDefault(t.Args, "{}"),
		State:        t.State,
		PollAfterSec: agentAssignmentPollSec,
	})
}

// buildTaskArgs describes the single container this task is.
func buildTaskArgs(db *gorm.DB, agent *models.Agent, run models.ClusterRun) (datatypes.JSON, error) {
	var cluster models.Cluster
	if err := db.Select("id", "docker_image", "docker_tag").
		Where("id = ?", agent.ClusterID).First(&cluster).Error; err != nil {
		return nil, err
	}

	var server models.Server
	if err := db.Select("id", "ssh_user").
		Where("id = ?", agent.ServerID).First(&server).Error; err != nil {
		return nil, err
	}

	// prepare_cluster.go pushes a cluster's assets to `$HOME/autoglue/clusters/
	// <id>` over SSH as Server.SSHUser, and runs the container with those paths
	// bound in. The agent hands mount strings straight to `docker run`, which
	// performs no shell expansion, so the control plane has to state the
	// expansion the remote shell used to do.
	home := "/root"
	if u := strings.TrimSpace(server.SSHUser); u != "" && u != "root" {
		home = "/home/" + u
	}
	clusterDir := home + "/autoglue/clusters/" + cluster.ID.String()

	args := agentTaskArgs{
		ClusterID: cluster.ID.String(),
		RunID:     run.ID.String(),
		Image:     cluster.DockerImage,
		Tag:       cluster.DockerTag,
		Target:    run.Action,
		// Deliberately empty. The SSH path `cd`s into the cluster directory only
		// so a relative bind mount resolves; docker's --workdir is container
		// side, and setting it to a host path would override the image's own
		// WORKDIR. The mounts below are absolute, so nothing needs it.
		WorkDir: "",
		Mounts: []string{
			home + "/.ssh:/root/.ssh",
			clusterDir + "/payload.json:/opt/gluekube/platform.json",
		},
	}

	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

// AgentTaskStart records that the agent has launched the container.
//
// It changes control-plane state, so it is a POST — but it releases nothing and
// proves nothing about the work. It exists so an operator can see a container
// id while a long make target is still running, and so the next assignment
// response can tell the agent which of the two in-flight states it is in.
func AgentTaskStart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, task, ok := agentTaskFromRequest(db, w, r)
		if !ok {
			return
		}

		var req agentTaskStartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid body")
			return
		}
		containerID := strings.TrimSpace(req.ContainerID)
		if containerID == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "container_id is required")
			return
		}

		res := db.Model(&models.AgentTask{}).
			Where("id = ? AND agent_id = ? AND state = ?",
				task.ID, agent.ID, models.AgentTaskStateAssigned).
			Updates(map[string]any{
				"state":        models.AgentTaskStateStarted,
				"container_id": containerID,
				"started_at":   time.Now().UTC(),
			})
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if res.RowsAffected == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Zero rows is re-checked rather than failed. The agent retries a lost
		// response, and a redelivery naming the same container is the same fact
		// arriving twice — only a *different* container, or a task that has
		// already ended, is a real disagreement worth surfacing.
		var current models.AgentTask
		if err := db.Where("id = ? AND agent_id = ?", task.ID, agent.ID).
			First(&current).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if current.State == models.AgentTaskStateStarted && current.ContainerID == containerID {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		utils.WriteError(w, http.StatusConflict, "conflict", "task is not startable in its current state")
	}
}

// AgentTaskLogs ingests a batch of container output.
//
// Chunks land in job_logs under subject_type='cluster_run' with the run's uuid,
// which is the whole point: readJobLogs, GetClusterRunLogs and the UI tail keep
// working with no change at all, and no reader can tell an agent wrote them
// rather than a River worker. Retention reaps them like any other chunk.
//
// Idempotency is a per-task watermark rather than a unique index on (task_id,
// seq). The agent posts serially in seq order per task, and given ordered
// delivery a monotonic watermark is exactly equivalent to a unique key — while
// costing no index on the highest-volume table in the schema and no dedupe
// table to sweep. Partial-batch recovery falls out of it: if the server writes
// 4..7 and dies before responding, the agent redelivers 4..10 and the watermark
// discards 4..7.
func AgentTaskLogs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, task, ok := agentTaskFromRequest(db, w, r)
		if !ok {
			return
		}

		var req agentTaskLogsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid body")
			return
		}
		if len(req.Chunks) > agentMaxLogChunks {
			utils.WriteError(w, http.StatusRequestEntityTooLarge, "too_many_chunks", "too many chunks in one batch")
			return
		}
		for _, c := range req.Chunks {
			// Caught explicitly rather than left to the global 10 MiB body cap,
			// which would surface as a truncated-JSON parse error that tells an
			// agent author nothing about what to fix.
			if len(c.Chunk) > agentMaxLogChunkBytes {
				utils.WriteError(w, http.StatusRequestEntityTooLarge, "chunk_too_large", "log chunk exceeds the maximum size")
				return
			}
		}

		accepted := task.LogSeq
		txErr := db.Transaction(func(tx *gorm.DB) error {
			// FOR UPDATE locks the watermark and re-authorizes in one statement,
			// so two concurrent batches cannot both read the same log_seq and
			// both decide their chunks are new.
			var locked models.AgentTask
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND agent_id = ?", task.ID, agent.ID).
				First(&locked).Error; err != nil {
				return err
			}
			accepted = locked.LogSeq

			rows := make([]models.JobLog, 0, len(req.Chunks))
			for _, c := range req.Chunks {
				// Compare against accepted, not locked.LogSeq. locked.LogSeq is
				// the pre-batch watermark and never moves inside this loop, so
				// filtering on it deduplicates against *previous* posts but not
				// against seqs already seen in this one — a batch carrying the
				// same seq twice would insert both, and the transcript readers
				// page by job_logs.id with no way to tell the copies apart.
				//
				// This also enforces ordering within the batch: a chunk whose
				// seq goes backwards is dropped rather than written out of
				// sequence, because the reader's cursor never revisits an id it
				// has already passed.
				if c.Seq <= accepted {
					continue
				}
				rows = append(rows, models.JobLog{
					// JobID 0 is the sanctioned "writer had no River job
					// context" value; an agent has none by construction.
					JobID:          0,
					OrganizationID: locked.OrganizationID,
					SubjectType:    models.JobLogSubjectClusterRun,
					SubjectID:      locked.RunID,
					Stream:         normalizeLogStream(c.Stream),
					Chunk:          c.Chunk,
				})
				if c.Seq > accepted {
					accepted = c.Seq
				}
			}
			if len(rows) == 0 {
				return nil
			}

			if err := tx.CreateInBatches(&rows, 200).Error; err != nil {
				return err
			}
			return tx.Model(&models.AgentTask{}).
				Where("id = ? AND log_seq < ?", locked.ID, accepted).
				UpdateColumn("log_seq", accepted).Error
		})
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "task not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Reporting the watermark rather than a count lets the agent advance its
		// own outbox precisely, including past chunks this call discarded.
		utils.WriteJSON(w, http.StatusOK, agentTaskLogsResponse{
			TaskID:             task.ID,
			AcceptedThroughSeq: accepted,
		})
	}
}

// normalizeLogStream keeps an unrecognised stream from reaching the readers.
// Docker gives the agent stdout and stderr; system is reserved for the control
// plane's own commentary, so anything unknown falls back to stdout, which is
// the value that renders as plain output.
func normalizeLogStream(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case agentLogStreamStderr:
		return agentLogStreamStderr
	case models.JobLogStreamSystem:
		return models.JobLogStreamSystem
	default:
		return models.JobLogStreamStdout
	}
}

// AgentTaskResult is the only thing that moves a task terminal, and therefore
// the only thing that releases the agent's slot.
//
// dead_lettered is a third outcome rather than a flavour of failure: it means
// nobody can say whether the work happened, which is not the same as knowing it
// did not. Without it, "never auto-retry" plus "only a terminal report frees the
// slot" would deadlock an agent permanently on the first ambiguous container.
// Nothing redrives it — recovery is a human posting a new run, which produces a
// fresh run and therefore a fresh task, because re-running a make target that
// may have half-applied is not something this system may assume is safe.
func AgentTaskResult(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, task, ok := agentTaskFromRequest(db, w, r)
		if !ok {
			return
		}

		var req agentTaskResultRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid body")
			return
		}

		state := strings.ToLower(strings.TrimSpace(req.State))
		switch state {
		case models.AgentTaskStateSucceeded,
			models.AgentTaskStateFailed,
			models.AgentTaskStateDeadLettered:
		default:
			utils.WriteError(w, http.StatusBadRequest, "bad_state", "state must be succeeded, failed or dead_lettered")
			return
		}

		now := time.Now().UTC()
		conflict := false

		txErr := db.Transaction(func(tx *gorm.DB) error {
			var locked models.AgentTask
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND agent_id = ?", task.ID, agent.ID).
				First(&locked).Error; err != nil {
				return err
			}

			if isTerminalTaskState(locked.State) {
				// Re-reporting the same outcome is a redelivery, which is the
				// normal cost of an at-least-once client. A *different* outcome
				// is two contradictory claims about one container, and picking
				// one would silently discard the other.
				if locked.State != state {
					conflict = true
				}
				return nil
			}

			updates := map[string]any{
				"state":     state,
				"error":     req.Error,
				"ended_at":  now,
				"exit_code": req.ExitCode,
			}
			if state == models.AgentTaskStateDeadLettered {
				updates["dead_letter_reason"] = models.AgentTaskDeadLetterAgentReported
			}
			if err := tx.Model(&models.AgentTask{}).
				Where("id = ?", locked.ID).
				Updates(updates).Error; err != nil {
				return err
			}

			// Propagate to the parent run in the same transaction, so a task
			// that is terminal and a run that is still running cannot be
			// observed together.
			if state == models.AgentTaskStateDeadLettered {
				return failRunFromDeadLetter(tx, locked.RunID,
					models.AgentTaskDeadLetterAgentReported, req.Error, now)
			}

			runStatus := clusterRunStatusSucceeded
			runError := ""
			if state == models.AgentTaskStateFailed {
				runStatus = models.ClusterRunStatusFailed
				runError = req.Error
			}
			return tx.Model(&models.ClusterRun{}).
				Where("id = ? AND status IN ?", locked.RunID,
					[]string{models.ClusterRunStatusQueued, models.ClusterRunStatusRunning}).
				Updates(map[string]any{
					"status":      runStatus,
					"error":       runError,
					"finished_at": now,
				}).Error
		})
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "task not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if conflict {
			utils.WriteError(w, http.StatusConflict, "conflict", "task already reported a different terminal state")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func isTerminalTaskState(s string) bool {
	switch s {
	case models.AgentTaskStateSucceeded,
		models.AgentTaskStateFailed,
		models.AgentTaskStateDeadLettered:
		return true
	}
	return false
}

// agentTaskFromRequest resolves {taskID} and proves it belongs to the calling
// agent.
//
// The ownership predicate is the authorization, and it is an AND on the same
// query rather than a load followed by a check: an agent that guesses another
// cluster's task id gets a 404, which is also the honest answer — that task does
// not exist as far as this principal is concerned.
func agentTaskFromRequest(db *gorm.DB, w http.ResponseWriter, r *http.Request) (*models.Agent, models.AgentTask, bool) {
	agent, ok := httpmiddleware.AgentFrom(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "agent required")
		return nil, models.AgentTask{}, false
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_task_id", "invalid task id")
		return nil, models.AgentTask{}, false
	}

	var task models.AgentTask
	if err := db.Where("id = ? AND agent_id = ?", taskID, agent.ID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteError(w, http.StatusNotFound, "not_found", "task not found")
			return nil, models.AgentTask{}, false
		}
		utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
		return nil, models.AgentTask{}, false
	}

	return agent, task, true
}
