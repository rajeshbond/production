package resourcemaster

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

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
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

	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Service.CreateResource(ctx, req, claims)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}

func (h *Handler) GetResourceByID(w http.ResponseWriter, r *http.Request) {
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

	resource, err := h.Service.GetResourceByID(ctx, id, claims.TenantID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, resource)
}

func (h *Handler) GetAllResources(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.Service.GetAllResources(ctx, claims.TenantID)
	if err != nil {
		response.Error(w, http.StatusFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list)
}

func (h *Handler) UpdateResorce(w http.ResponseWriter, r *http.Request) {
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
	var req UpdateResourceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrInvalidRequest.Error())
		return
	}

	if err := h.Service.UpdateResorce(ctx, req, claims); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "updated sucessfully")

}

func (h *Handler) DeleteResource(w http.ResponseWriter, r *http.Request) {
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
	if err := h.Service.DeleteResource(ctx, id, claims); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "Deleted Sucessfully")

}
