package api

import (
	httpPprof "net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func mountPprofRoutes(r chi.Router) {
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
}
