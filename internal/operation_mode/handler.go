package operationmode

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	var req CreateOperationModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// tenantID := int64(1) // from auth
	// userID := int64(1)

	if err := h.Service.Create(ctx, req, claims); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// w.WriteHeader(http.StatusCreated,oK)
	response.JSON(w, http.StatusCreated, "Created")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	var req UpdateOperationModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tenantID := int64(1)
	userID := int64(1)

	if err := h.Service.Update(r.Context(), tenantID, userID, req); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	tenantID := int64(1)
	userID := int64(1)

	if err := h.Service.Delete(r.Context(), tenantID, userID, id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
