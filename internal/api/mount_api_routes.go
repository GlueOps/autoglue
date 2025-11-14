package api

import (
	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountAPIRoutes(r chi.Router, db *gorm.DB, jobs *bg.Jobs) {
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			authUser := httpmiddleware.AuthMiddleware(db, false)
			authOrg := httpmiddleware.AuthMiddleware(db, true)

			// shared basics
			mountMetaRoutes(v1)
			mountAuthRoutes(v1, db)

			// admin
			mountAdminRoutes(v1, db, jobs, authUser)

			// user/org scoped
			mountMeRoutes(v1, db, authUser)
			mountOrgRoutes(v1, db, authUser, authOrg)

			mountCredentialRoutes(v1, db, authOrg)
			mountSSHRoutes(v1, db, authOrg)
			mountServerRoutes(v1, db, authOrg)
			mountTaintRoutes(v1, db, authOrg)
			mountLabelRoutes(v1, db, authOrg)
			mountAnnotationRoutes(v1, db, authOrg)
			mountNodePoolRoutes(v1, db, authOrg)
			mountDNSRoutes(v1, db, authOrg)
		})
	})
}
