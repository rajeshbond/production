package mould

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	Service *Service

	// tokenAuth *jwtauth.JWTAuth
}

// func NewHandler(tokenAuth *jwtauth.JWTAuth, service *Service) *Handler {
// 	return &Handler{tokenAuth: tokenAuth, Service: service}
// }

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateMold(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}

	var req CreateMoldRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.ToLower(req.Type) != "mold" {
		http.Error(w, "invalid request body: type must be Mold", http.StatusBadRequest)
		return
	}

	if req.MoldNo == "" {
		response.Error(w, http.StatusBadRequest, "Mould no required")
		return
	}

	id, err := h.Service.Create(ctx, &req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)
}
