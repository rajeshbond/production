package resourcetype

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		Service:   service}
}

func (h *Handler) CreateResourceType(w http.ResponseWriter, r *http.Request) {
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
	var req CreateResourceTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.Service.CreateResourceType(ctx, req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}

func (h *Handler) UpdateResourceType(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateResourceTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.InvalidRequestBody)
		return
	}

	response.JSON(w, http.StatusOK, req)

}

func (h *Handler) DeleteResourceType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, auth.UnAuthorised)
		return
	}

	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, auth.UnAuthorised)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "Id required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	deletedID, err := h.Service.DeleteResourceType(ctx, id, claims)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, deletedID)

}
