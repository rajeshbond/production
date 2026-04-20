package operationdefectmap

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

func (h *Handler) CreateOperationWithDefect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// fmt.Println("Rajesh Bondgilwar")

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	var req OperationDefectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.InvalidRequest)
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	// if claims.Role != "tenantadmin" {
	// 	response.Error(w, http.StatusForbidden, "only tenantadmin allowed")
	// 	return
	// }

	resp, err := h.Service.CreateOperationWithDefect(ctx, req, *claims)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}
