package tools

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateTools(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {

		response.Error(w, http.StatusForbidden, auth.UnAuthorised)
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusForbidden, auth.UnAuthorised)
		return
	}

	var req CreateToolRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.ToLower(req.Type) != "tool" {
		http.Error(w, "invalid request body: type must be tool", http.StatusBadRequest)
		return
	}

	// 4. Call service (transaction handled inside service)
	id, err := h.Service.Create(ctx, &req, claims)
	if err != nil {
		http.Error(w, "failed to create tool: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, id)

}
