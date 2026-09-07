package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func InitializeDebugRoutes(router *chi.Mux) {
	if router == nil {
		return
	}
	debugRouter := chi.NewMux()
	debugRouter.Get("/pprof/", http.HandlerFunc(pprof.Index))
	debugRouter.Get("/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	debugRouter.Get("/pprof/profile", http.HandlerFunc(pprof.Profile))
	debugRouter.Get("/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	debugRouter.Get("/pprof/trace", http.HandlerFunc(pprof.Trace))
	router.Mount("/debug/pprof", debugRouter)
}
