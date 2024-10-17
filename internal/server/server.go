package server

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
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

	r.Route("/v1", func(r chi.Router) {
		r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/pprof/mutex", pprof.Handler("mutex"))
		r.Handle("/pprof/heap", pprof.Handler("heap"))
		r.Handle("/pprof/block", pprof.Handler("block"))
		r.Handle("/pprof/allocs", pprof.Handler("allocs"))

		r.Route("/{mount}", func(r chi.Router) {
			r.Get("/config", h.GetKVConfig)
			r.Post("/config", h.UpdateKVConfig)

			r.Get("/data/{path}", h.GetKVSecret)
			r.Post("/data/{path}", h.CreateKVSecret)
			r.Put("/data/{path}", h.UpdateKVSecret)
			r.Delete("/data/{path}", h.DeleteLatestKVSecret)

			r.Post("/delete/{path}", h.DeleteKVSecret)
			r.Post("/destroy/{path}", h.DestroyKVSecret)

			r.Get("/metadata/{path}", h.GetKVMetadata)
			r.Post("/metadata/{path}", h.UpdateKVMetadata)
			r.Delete("/metadata/{path}", h.DeleteKVMetadata)

			r.Get("/subkeys/{path}", h.GetKVSubkeys)
			r.Post("/subkeys/{path}", h.CreateKVSubkeys)
		})

		r.Route("/auth/token", func(r chi.Router) {
			r.Get("/accessors/", h.GetTokenAccessors)
			r.Post("/create", h.CreateToken)
			r.Post("/create-orphan", h.CreateOrphanToken)
			r.Post("/create/{role_name}", h.CreateRoleToken)
			r.Get("/lookup", h.LookupToken)
			r.Post("/lookup", nil)
			r.Post("/lookup-accessor", nil)
			r.Get("/lookup-self", h.LookupToken)
			r.Post("/lookup-self", h.LookupSelfToken)
			r.Post("/renew", h.RenewToken)
			r.Post("/renew-accessor", h.RenewAccessorToken)
			r.Post("/renew-self", h.RenewSelfToken)
			r.Post("/revoke", h.RevokeToken)
			r.Post("/revoke-accessor", h.RevokeAccessorToken)
			r.Post("/revoke-orphan", h.RevokeOrphanToken)
			r.Post("/revoke-self", h.RevokeSelfToken)
			r.Get("/roles/", h.GetRolesToken)
			r.Get("/roles/{role_name}", h.GetRoleByNameToken)
			r.Post("/roles/{role_name}", h.CreateRoleByNameToken)
			r.Delete("/roles/{role_name}", h.DeleteRoleByNameToken)
			r.Post("/tidy", h.TidyToken)
		})

		r.Route("/sys/tools", func(r chi.Router) {
			r.Post("/hash", nil)
			r.Post("/hash/{urlalgorithm}", nil)
			r.Post("/random", nil)
			r.Post("/random/{source}", nil)
			r.Post("/random/{source}/{urlbytes}", nil)
			r.Post("/random/{urlbytes}", nil)
		})

		r.Route("/sys", func(r chi.Router) {
			r.Get("/mounts", h.GetMounts)
			r.Get("/mounts/{path}", h.GetMount)
			r.Post("/mounts/{path}", h.CreateMount)
			r.Delete("/mounts/{path}", h.DeleteMount)

			r.Post("/seal", h.Seal)
			r.Get("/seal-status", h.SealStatus)
			r.Post("/unseal", h.Unseal)
			r.Post("/init", h.Init)
			r.Get("/health", h.Health)

			r.Get("/metrics", promhttp.Handler().ServeHTTP)
			r.HandleFunc("/pprof/*", pprof.Index)
			r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
			r.HandleFunc("/pprof/profile", pprof.Profile)
			r.HandleFunc("/pprof/symbol", pprof.Symbol)
			r.HandleFunc("/pprof/trace", pprof.Trace)
		})
	})

	srv.server.Handler = r

	return srv
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
