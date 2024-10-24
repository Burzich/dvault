package server

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server  http.Server
	handler DVaultHandler
}

func NewServer(addr string, h DVaultHandler) *Server {
	srv := &Server{
		server: http.Server{
			Addr: addr,
		},
		handler: h,
	}

	r := chi.NewMux()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/v1/sys/seal-status", h.SealStatus)
	r.Post("/v1/sys/unseal", h.Unseal)
	r.Post("/v1/sys/init", h.Init)

	r.Route("/v1", func(r chi.Router) {
		r.Use(h.AuthMiddleware())
		r.Route("/{mount}", func(r chi.Router) {
			r.Get("/config", h.GetKVConfig)
			r.Post("/config", h.UpdateKVConfig)

			r.Get("/data/{path}", h.GetKVSecret)
			r.Post("/data/{path}", h.CreateKVSecret)
			r.Delete("/data/{path}", h.DeleteLatestKVSecret)

			r.Post("/delete/{path}", h.DeleteKVSecret)
			r.Post("/destroy/{path}", h.DestroyKVSecret)

			r.Get("/metadata/{path}", h.GetKVMetadata)
			r.Post("/metadata/{path}", h.UpdateKVMetadata)
			r.Delete("/metadata/{path}", h.DeleteKVMetadata)

			r.Get("/subkeys/{path}", h.GetKVSubkeys)
			r.Post("/subkeys/{path}", h.CreateKVSubkeys)
		})

		r.Route("/sys", func(r chi.Router) {
			r.Get("/mounts", h.GetMounts)
			r.Get("/mounts/{path}", h.GetMount)
			r.Post("/mounts/{path}", h.CreateMount)
			r.Delete("/mounts/{path}", h.DeleteMount)

			r.Post("/seal", h.Seal)
			r.Get("/health", h.Health)

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
		})

	})

	srv.server.Handler = r

	return srv
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return s.server.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
