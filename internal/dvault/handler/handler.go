package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Burzich/dvault/internal/dvault"
	"github.com/Burzich/dvault/internal/dvault/kv"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	dVault *dvault.DVault
}

func NewHandler(dVault *dvault.DVault) Handler {
	return Handler{dVault: dVault}
}

func (h Handler) GetKVConfig(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")

	response, err := h.dVault.GetKVConfig(r.Context(), mount)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) UpdateKVConfig(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")

	var updateConfig UpdateKVConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&updateConfig); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := h.dVault.UpdateKVConfig(r.Context(), mount, kv.Config{
		CasRequired:        updateConfig.CasRequired,
		DeleteVersionAfter: updateConfig.DeleteVersionAfter,
		MaxVersions:        updateConfig.MaxVersions,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetKVSecret(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")
	version := r.URL.Query().Get("version")

	var response dvault.Response
	if version != "" {
		v, err := strconv.Atoi(version)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err = h.dVault.GetKVSecretByVersion(r.Context(), mount, secretPath, v)
		if err != nil {
			h.handleError(w, r, err)
			return
		}
	} else {
		var err error
		response, err = h.dVault.GetKVSecret(r.Context(), mount, secretPath)
		if err != nil {
			h.handleError(w, r, err)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreateKVSecret(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	var createKV CreateKVSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&createKV); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := h.dVault.SaveKVSecret(r.Context(), mount, secretPath, createKV.Data, createKV.Options.CAS)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteLatestKVSecret(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	response, err := h.dVault.DeleteKVSecret(r.Context(), mount, secretPath)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteKVSecret(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	var deleteKVSecret DeleteKVSecret
	if err := json.NewDecoder(r.Body).Decode(&deleteKVSecret); err != nil {
		h.handleError(w, r, err)
		return
	}

	response, err := h.dVault.DeleteKVSecretByVersion(r.Context(), mount, secretPath, deleteKVSecret.Versions)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DestroyKVSecret(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	var destroyKVSecret DestroyKVSecret
	if err := json.NewDecoder(r.Body).Decode(&destroyKVSecret); err != nil {
		h.handleError(w, r, err)
		return
	}

	response, err := h.dVault.DestroyKVSecret(r.Context(), mount, secretPath, destroyKVSecret.Versions)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetKVMetadata(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	response, err := h.dVault.GetKVMeta(r.Context(), mount, secretPath)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) UpdateKVMetadata(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	var updateKVMetadata UpdateKVMetadata
	if err := json.NewDecoder(r.Body).Decode(&updateKVMetadata); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := h.dVault.UpdateKVMeta(r.Context(), mount, secretPath, kv.Meta{
		CasRequired:        updateKVMetadata.CasRequired,
		DeleteVersionAfter: updateKVMetadata.DeleteVersionAfter,
		CustomMetadata:     updateKVMetadata.CustomMetadata,
		MaxVersions:        updateKVMetadata.MaxVersions,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteKVMetadata(w http.ResponseWriter, r *http.Request) {
	mount := chi.URLParam(r, "mount")
	secretPath := chi.URLParam(r, "path")

	response, err := h.dVault.DeleteKVMeta(r.Context(), mount, secretPath)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetKVSubkeys(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateKVSubkeys(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) GetTokenAccessors(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateOrphanToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateRoleToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) LookupToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) LookupSelfToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RenewToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RenewAccessorToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RenewSelfToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RevokeAccessorToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RevokeOrphanToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) RevokeSelfToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) GetRolesToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) GetRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) DeleteRoleByNameToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) TidyToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) Unseal(w http.ResponseWriter, r *http.Request) {
	var unsealRequest UnsealRequest
	if err := json.NewDecoder(r.Body).Decode(&unsealRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.dVault.Unseal(r.Context(), dvault.Unseal{
		Key:     unsealRequest.Key,
		Migrate: unsealRequest.Migrate,
		Reset:   unsealRequest.Reset,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) Seal(w http.ResponseWriter, r *http.Request) {
	response, err := h.dVault.Seal(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) SealStatus(w http.ResponseWriter, r *http.Request) {
	sealStatus, err := h.dVault.SealStatus(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(sealStatus); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) Init(w http.ResponseWriter, r *http.Request) {
	var initRequest InitRequest
	if err := json.NewDecoder(r.Body).Decode(&initRequest); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.dVault.Init(r.Context(), dvault.Init{
		PgpKeys:           initRequest.PgpKeys,
		RecoveryPgpKeys:   initRequest.RecoveryPgpKeys,
		RecoveryShares:    initRequest.RecoveryShares,
		RecoveryThreshold: initRequest.RecoveryThreshold,
		RootTokenPgpKey:   initRequest.RootTokenPgpKey,
		SecretShares:      initRequest.SecretShares,
		SecretThreshold:   initRequest.SecretThreshold,
		StoredShares:      initRequest.StoredShares,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetMounts(w http.ResponseWriter, r *http.Request) {
	_, err := h.dVault.SealStatus(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var sealStatus SealStatusResponse
	if err := json.NewEncoder(w).Encode(sealStatus); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetMount(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) CreateMount(w http.ResponseWriter, r *http.Request) {
	secretPath := chi.URLParam(r, "path")

	var createMount CreateMount
	if err := json.NewDecoder(r.Body).Decode(&createMount); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := h.dVault.CreateMount(r.Context(), secretPath, dvault.CreateMount{
		Config:                createMount.Config,
		Description:           createMount.Description,
		ExternalEntropyAccess: createMount.ExternalEntropyAccess,
		Local:                 createMount.Local,
		Options:               createMount.Options,
		PluginName:            createMount.PluginName,
		PluginVersion:         createMount.PluginVersion,
		SealWrap:              createMount.SealWrap,
		Type:                  createMount.Type,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteMount(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h Handler) handleError(rw http.ResponseWriter, r *http.Request, err error) {
	var b *strconv.NumError

	switch {
	case errors.As(err, &b):
		rw.WriteHeader(http.StatusBadRequest)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}

	errorResponse := Error{[]string{err.Error()}}
	if err = json.NewEncoder(rw).Encode(errorResponse); err != nil {
		return
	}
}
