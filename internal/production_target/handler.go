package productiontarget

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
		Service:   service,
	}
}

func (h *Handler) CreateProductionTarget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}
	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	var req CreateProductionTargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateProductionTarget(ctx, &req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.InvalidRequestBody)
		return
	}

	response.JSON(w, http.StatusCreated, id)

}

func (h *Handler) UpdateProductionTarget(w http.ResponseWriter, r *http.Request) {

	var req UpdateProductionTargetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}
	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	err := h.Service.UpdateProductionTarget(
		r.Context(),
		&req,
		claims,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("updated successfully"))
}

func (h *Handler) GetProductionTargetByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	tenantID := r.Context().Value("tenant_id").(int64)

	data, err := h.Service.GetProductionTargetByID(
		r.Context(),
		tenantID,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetAllProductionTargets(w http.ResponseWriter, r *http.Request) {

	tenantID := r.Context().Value("tenant_id").(int64)

	data, err := h.Service.GetAllProductionTargets(
		r.Context(),
		tenantID,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}
