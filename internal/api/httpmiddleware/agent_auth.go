package httpmiddleware

import (
	"net/http"
	"time"

	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Agent credential headers. Three of them because the credential is a tuple:
// see auth.AgentCredential for why no part of it is derivable from the others.
const (
	agentIDHeader     = "X-Agent-ID"
	agentKeyHeader    = "X-Agent-KEY"
	agentSecretHeader = "X-Agent-SECRET"
)

// AgentAuth authenticates a bastion agent and resolves it to a cluster-scoped
// principal.
//
// It is deliberately a separate middleware rather than a fourth branch of
// AuthMiddleware. That middleware's org resolution is entirely user-membership
// reasoning — X-Org-ID, the {id} URL param, the single-membership fallback —
// and none of it has an agent analogue. More importantly an agent must never
// resolve to an org-wide principal: reaching exactly one cluster is the whole
// point of the credential, and a leaked agent that could act on the org would
// be no better than the org-wide bastion key this replaces.
//
// The context gets the agent and its organization, and deliberately neither a
// user nor roles. That makes RequireAuthenticatedUser, RequirePlatformAdmin and
// RequireRole all reject an agent structurally, should one of them ever be
// composed onto an agent route by accident.
func AgentAuth(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			agent := auth.ValidateAgentKeyPair(
				r.Header.Get(agentIDHeader),
				r.Header.Get(agentKeyHeader),
				r.Header.Get(agentSecretHeader),
				db,
			)
			if agent == nil {
				// One response for every failure — missing header, unknown key,
				// wrong secret, mismatched id, revoked, expired. Distinguishing
				// them would tell a caller holding part of a credential which
				// part still works.
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid agent credentials")
				return
			}

			// Org is loaded rather than stubbed from agent.OrganizationID so
			// that OrgFrom yields a real row: a half-populated organization in
			// the context is the kind of thing a shared helper reads a name off
			// months later.
			var org models.Organization
			if err := db.First(&org, "id = ?", agent.OrganizationID).Error; err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid agent credentials")
				return
			}

			touchAgentLastSeen(db, agent)

			ctx := WithAgent(r.Context(), agent)
			ctx = WithOrg(ctx, &org)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// touchAgentLastSeen records liveness as a side effect of authenticating.
//
// This is the honest heartbeat precisely because the agent did not ask for it:
// every loop it runs authenticates, so the timestamp is control-plane observed
// contact rather than a claim the agent makes about itself. The reconcile
// report updates the same field, but only on the config loop's timer — an agent
// that is polling for work while its config loop is wedged still looks alive
// here, which is the truth.
//
// UpdateColumn rather than Update so updated_at keeps meaning "the control
// plane changed this row" instead of "a bastion polled".
//
// A failed write is not a failed request. The agent's work does not depend on
// the control plane remembering when it last called.
func touchAgentLastSeen(db *gorm.DB, agent *models.Agent) {
	now := time.Now()
	if err := db.Model(&models.Agent{}).
		Where("id = ?", agent.ID).
		UpdateColumn("last_seen_at", now).Error; err != nil {
		log.Warn().Err(err).Str("agent_id", agent.ID.String()).Msg("agent last_seen_at update failed")
		return
	}
	agent.LastSeenAt = &now
}
