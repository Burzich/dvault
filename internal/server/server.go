package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server http.Server
}

func NewServer(addr string) *Server {
	srv := &Server{
		server: http.Server{
			Addr: addr,
		},
	}

	r := chi.NewMux()
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	r.Route("/kv", func(r chi.Router) {
		r.Get("/config", nil)
		r.Post("/config", nil)

		r.Get("/data/{path}", nil)
		r.Post("/data/{path}", nil)
		r.Delete("/data/{path}", nil)

		r.Post("/delete/{key}", nil)
		r.Post("/destroy/{key}", nil)

		r.Get("/metadata/{path}", nil)
		r.Post("/metadata/{path}", nil)
		r.Delete("/metadata/{path}", nil)
		r.Get("/metadata/{path}/", nil)

		r.Get("/subkeys/{path}", nil)
		r.Post("/subkeys/{path}", nil)
	})

	srv.server.Handler = r

	return srv
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
