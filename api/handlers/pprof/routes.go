package pprof

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
)

func RegisterPprofRoutes(router *mux.Router) {
	// @Summary Pprof Index
	// @Tags Internal, Debug
	// @Router /debug/pprof/ [get]
	router.HandleFunc("/debug/pprof/", pprof.Index)

	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// @Summary CPU Profile
	// @Tags Internal, Debug
	// @Router /debug/pprof/profile [get]
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	// Manually add support for paths linked to by index page at /debug/pprof/
	// @Summary Goroutines
	// @Tags Internal, Debug
	// @Router /debug/pprof/goroutine [get]
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))

	// @Summary      Heap profiling
	// @Tags         Internal, Debug
	// @Router       /debug/pprof/heap [get]
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
}
