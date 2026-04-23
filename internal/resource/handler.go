package resource

import (
	"encoding/json"
	"net/http"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
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

	var req CreateResourceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if req.ResourceCode == "" || req.ResourceType == "" {
		response.BadRequest(w, "resource_code and resource_type required")
		return
	}

	id, err := h.Service.Create(ctx, &req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}
