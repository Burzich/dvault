package server

import (
	"net/http"
)

type DVaultHandler interface {
	GetKVConfig(w http.ResponseWriter, r *http.Request)
	CreateKVConfig(w http.ResponseWriter, r *http.Request)

	GetKVSecret(w http.ResponseWriter, r *http.Request)
	CreateKVSecret(w http.ResponseWriter, r *http.Request)
	DeleteKVSecret(w http.ResponseWriter, r *http.Request)

	DeleteKV(w http.ResponseWriter, r *http.Request)
	DestroyKV(w http.ResponseWriter, r *http.Request)

	GetKVMetadata(w http.ResponseWriter, r *http.Request)
	CreateKVMetadata(w http.ResponseWriter, r *http.Request)
	DeleteKVMetadata(w http.ResponseWriter, r *http.Request)
	UpdateKVMetadata(w http.ResponseWriter, r *http.Request)

	GetKVSubkeys(w http.ResponseWriter, r *http.Request)
	CreateKVSubkeys(w http.ResponseWriter, r *http.Request)
	/*
		r.Route("/auth/token", func(r chi.Router) {
			r.Get("accessors/", nil)
			r.Post("/create", nil)
			r.Post("/create-orphan", nil)
			r.Post("/create/{role_name}", nil)
			r.Get("/lookup", nil)
			r.Post("/lookup", nil)
			r.Post("/lookup-accessor", nil)
			r.Get("/lookup-self", nil)
			r.Post("/lookup-self", nil)
			r.Post("/renew", nil)
			r.Post("/renew-accessor", nil)
			r.Post("/renew-self", nil)
			r.Post("/revoke", nil)
			r.Post("/revoke-accessor", nil)
			r.Post("/revoke-orphan", nil)
			r.Post("/revoke-self", nil)
			r.Get("/roles/", nil)
			r.Get("/roles/{role_name}", nil)
			r.Post("/roles/{role_name}", nil)
			r.Delete("/roles/{role_name}", nil)
			r.Post("/tidy", nil)
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
			r.Post("/seal", nil)
			r.Get("/seal-status", nil)
			r.Post("/unseal", nil)
			r.Get("/health", nil)
		})
	*/
}
