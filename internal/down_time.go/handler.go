package downtime

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
	"github.com/rajesh_bond/production/internal/common/utils"
)

type Handler struct {
	Service   *Service
	tokenAuth *jwtauth.JWTAuth
}

// Struct connstructor
func NewHandler(service *Service, tokenAuth *jwtauth.JWTAuth) *Handler {
	return &Handler{
		Service:   service,
		tokenAuth: tokenAuth,
	}
}

func (h *Handler) CreateDownTime(w http.ResponseWriter, r *http.Request) {
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

	var req BulkCreateDownTimeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}
	if err := utils.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call service

	result, err := h.Service.CreateDownTime(ctx, req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, result)

}
