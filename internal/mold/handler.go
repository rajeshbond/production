package mold

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	Service   *Service
	tokenAuth *jwtauth.JWTAuth
}

func NewHandler(service *Service, tokenAuth *jwtauth.JWTAuth) *Handler {
	return &Handler{
		tokenAuth: tokenAuth,
		Service:   service,
	}
}

func (h *Handler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}
	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, auth.ErrUnauthorized.Error())
		return
	}

	var req BulkCreateMoldRequest

	inserted, skipped, err := h.Service.BulkCreate(ctx, req, claims)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]interface{}{
		"inserted": inserted,
		"skipped":  skipped,
	}

	response.JSON(w, http.StatusCreated, resp)

}
