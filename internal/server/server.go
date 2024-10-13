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

func NewServer(addr string, handler DVaultHandler) *Server {
	srv := &Server{
		server: http.Server{
			Addr: addr,
		},
		handler: handler,
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
		r.Get("/config", handler.GetKVConfig)
		r.Post("/config", handler.CreateKVConfig)

		r.Get("/data/{path}", handler.GetKVSecret)
		r.Post("/data/{path}", handler.CreateKVSecret)
		r.Delete("/data/{path}", handler.DeleteKVSecret)

		r.Post("/delete/{key}", handler.DeleteKV)
		r.Post("/destroy/{key}", handler.DestroyKV)

		r.Get("/metadata/{path}", handler.GetKVMetadata)
		r.Post("/metadata/{path}", handler.CreateKVMetadata)
		r.Delete("/metadata/{path}", handler.DeleteKVMetadata)
		r.Get("/metadata/{path}/", nil)

		r.Get("/subkeys/{path}", handler.GetKVSubkeys)
		r.Post("/subkeys/{path}", handler.CreateKVSubkeys)
	})

	r.Route("/auth/token", func(r chi.Router) {
		r.Get("accessors/", handler.GetTokenAccessors)
		r.Post("/create", handler.CreateToken)
		r.Post("/create-orphan", handler.CreateOrphanToken)
		r.Post("/create/{role_name}", handler.CreateRoleToken)
		r.Get("/lookup", handler.LookupToken)
		r.Post("/lookup", nil)
		r.Post("/lookup-accessor", nil)
		r.Get("/lookup-self", handler.LookupToken)
		r.Post("/lookup-self", handler.LookupSelfToken)
		r.Post("/renew", handler.RenewToken)
		r.Post("/renew-accessor", handler.RenewAccessorToken)
		r.Post("/renew-self", handler.RenewSelfToken)
		r.Post("/revoke", handler.RevokeToken)
		r.Post("/revoke-accessor", handler.RevokeAccessorToken)
		r.Post("/revoke-orphan", handler.RevokeOrphanToken)
		r.Post("/revoke-self", handler.RevokeSelfToken)
		r.Get("/roles/", handler.GetRolesToken)
		r.Get("/roles/{role_name}", handler.GetRoleByNameToken)
		r.Post("/roles/{role_name}", handler.CreateRoleByNameToken)
		r.Delete("/roles/{role_name}", handler.DeleteRoleByNameToken)
		r.Post("/tidy", handler.TidyToken)
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
		r.Post("/seal", handler.Seal)
		r.Get("/seal-status", handler.SealStatus)
		r.Post("/unseal", handler.Unseal)
		r.Get("/health", handler.Health)
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
