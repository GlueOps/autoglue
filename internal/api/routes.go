package api

import (
	httpPprof "net/http/pprof"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/handlers/authn"
	"github.com/glueops/autoglue/internal/handlers/health"
	"github.com/glueops/autoglue/internal/handlers/orgs"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/ui"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/api", func(api chi.Router) {
		api.Get("/healthz", health.Check)

		api.Route("/v1", func(v1 chi.Router) {
			secret := viper.GetString("authentication.jwt_secret")
			authMW := middleware.AuthMiddleware(secret)

			v1.Route("/auth", func(a chi.Router) {
				a.Post("/login", authn.Login)
				a.Post("/register", authn.Register)
				a.Post("/introspect", authn.Introspect)
				a.Post("/password/forgot", authn.RequestPasswordReset)
				a.Post("/password/reset", authn.ConfirmPasswordReset)
				a.Get("/verify", authn.VerifyEmail)
				a.Post("/verify/resend", authn.ResendVerification)

				a.Group(func(pr chi.Router) {
					pr.Use(authMW)
					pr.Post("/refresh", authn.Refresh)
					pr.Post("/logout", authn.Logout)
					pr.Post("/logout_all", authn.LogoutAll)
					pr.Get("/me", authn.Me)
					pr.Post("/password/change", authn.ChangePassword)
					pr.Post("/refresh/rotate", authn.RotateRefreshToken)
				})
			})

			v1.Route("/orgs", func(o chi.Router) {
				o.Use(authMW)
				o.Post("/", orgs.CreateOrganization)
				o.Get("/", orgs.ListOrganizations)
			})
		})
	})

	r.Route("/debug/pprof", func(pr chi.Router) {
		pr.Get("/", httpPprof.Index)
		pr.Get("/cmdline", httpPprof.Cmdline)
		pr.Get("/profile", httpPprof.Profile)
		pr.Get("/symbol", httpPprof.Symbol)
		pr.Get("/trace", httpPprof.Trace)

		pr.Handle("/allocs", httpPprof.Handler("allocs"))
		pr.Handle("/block", httpPprof.Handler("block"))
		pr.Handle("/goroutine", httpPprof.Handler("goroutine"))
		pr.Handle("/heap", httpPprof.Handler("heap"))
		pr.Handle("/mutex", httpPprof.Handler("mutex"))
		pr.Handle("/threadcreate", httpPprof.Handler("threadcreate"))
	})

	if config.IsUIDev() {
		if h, err := ui.DevProxy("http://localhost:5173"); err == nil {
			r.NotFound(h.ServeHTTP)
		}
	} else {
		if h, err := ui.SPAHandler(); err == nil {
			r.NotFound(h.ServeHTTP)
		}
	}
}
