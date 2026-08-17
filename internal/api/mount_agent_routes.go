package api

import (
	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// mountAgentRoutes mounts the bastion agent plane.
//
// It takes its own auth rather than the authOrg the other mounts share. An
// agent is a cluster-scoped principal and must never resolve to an org-wide
// one — that is the entire reason this exists, since the credential it replaces
// is an org key with 24 hours of life sitting on a bastion.
//
// Enrolment sits outside the authenticated group because it is what issues the
// credential the group checks. Its ticket is the credential.
//
// Nothing here is annotated for swag, deliberately: no human interacts with
// these endpoints, and an operation in docs/openapi.json is an invitation to.
func mountAgentRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/agent", func(a chi.Router) {
		a.Post("/enroll", handlers.EnrollAgent(db))

		a.Group(func(auth chi.Router) {
			auth.Use(httpmiddleware.AgentAuth(db))

			// GET where the mutation lands on the bastion's own disk rather
			// than on control-plane state; POST for everything that changes
			// what the control plane believes.
			auth.Get("/sync", handlers.AgentSync(db))
			auth.Get("/assignment", handlers.AgentAssignment(db))

			auth.Post("/reconcile-report", handlers.AgentReconcileReport(db))
			auth.Post("/tasks/{taskID}/start", handlers.AgentTaskStart(db))
			auth.Post("/tasks/{taskID}/logs", handlers.AgentTaskLogs(db))
			auth.Post("/tasks/{taskID}/result", handlers.AgentTaskResult(db))
		})
	})
}
