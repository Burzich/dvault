package handler

import (
	"net/http"

	"github.com/Burzich/dvault/internal/dvault"
)

type Handler struct {
	dVault *dvault.DVault
}

func (h Handler) GetKVConfig(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateKVConfig(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetKVSecret(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateKVSecret(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) DeleteKVSecret(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) DeleteKV(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) DestroyKV(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetKVMetadata(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateKVMetadata(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) DeleteKVMetadata(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) UpdateKVMetadata(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetKVSubkeys(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateKVSubkeys(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetTokenAccessors(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateOrphanToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateRoleToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) LookupToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) LookupSelfToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RenewToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RenewAccessorToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RenewSelfToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RevokeAccessorToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RevokeOrphanToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) RevokeSelfToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetRolesToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) GetRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) CreateRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) DeleteRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) TidyToken(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) Unseal(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) Seal(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) SealStatus(w http.ResponseWriter, r *http.Request) {
	return
}

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	return
}

func NewHandler(dvault *dvault.DVault) Handler {
	return Handler{dVault: dvault}
}
