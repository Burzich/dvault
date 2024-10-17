package server

import (
	"net/http"
)

type DVaultHandler interface {
	GetKVConfig(w http.ResponseWriter, r *http.Request)
	UpdateKVConfig(w http.ResponseWriter, r *http.Request)
	GetKVSecret(w http.ResponseWriter, r *http.Request)
	CreateKVSecret(w http.ResponseWriter, r *http.Request)
	UpdateKVSecret(w http.ResponseWriter, r *http.Request)
	DeleteLatestKVSecret(w http.ResponseWriter, r *http.Request)
	DeleteKVSecret(w http.ResponseWriter, r *http.Request)
	DestroyKVSecret(w http.ResponseWriter, r *http.Request)
	GetKVMetadata(w http.ResponseWriter, r *http.Request)
	UpdateKVMetadata(w http.ResponseWriter, r *http.Request)
	DeleteKVMetadata(w http.ResponseWriter, r *http.Request)
	GetKVSubkeys(w http.ResponseWriter, r *http.Request)
	CreateKVSubkeys(w http.ResponseWriter, r *http.Request)

	GetMounts(w http.ResponseWriter, r *http.Request)
	GetMount(w http.ResponseWriter, r *http.Request)
	CreateMount(w http.ResponseWriter, r *http.Request)
	DeleteMount(w http.ResponseWriter, r *http.Request)

	GetTokenAccessors(w http.ResponseWriter, r *http.Request)
	CreateToken(w http.ResponseWriter, r *http.Request)
	CreateOrphanToken(w http.ResponseWriter, r *http.Request)
	CreateRoleToken(w http.ResponseWriter, r *http.Request)
	LookupToken(w http.ResponseWriter, r *http.Request)
	LookupSelfToken(w http.ResponseWriter, r *http.Request)
	RenewToken(w http.ResponseWriter, r *http.Request)
	RenewAccessorToken(w http.ResponseWriter, r *http.Request)
	RenewSelfToken(w http.ResponseWriter, r *http.Request)
	RevokeToken(w http.ResponseWriter, r *http.Request)
	RevokeAccessorToken(w http.ResponseWriter, r *http.Request)
	RevokeOrphanToken(w http.ResponseWriter, r *http.Request)
	RevokeSelfToken(w http.ResponseWriter, r *http.Request)
	GetRolesToken(w http.ResponseWriter, r *http.Request)
	GetRoleByNameToken(w http.ResponseWriter, r *http.Request)
	CreateRoleByNameToken(w http.ResponseWriter, r *http.Request)
	DeleteRoleByNameToken(w http.ResponseWriter, r *http.Request)
	TidyToken(w http.ResponseWriter, r *http.Request)

	Unseal(w http.ResponseWriter, r *http.Request)
	Seal(w http.ResponseWriter, r *http.Request)
	SealStatus(w http.ResponseWriter, r *http.Request)
	Init(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
}
