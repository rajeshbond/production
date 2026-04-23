package fixture

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

func (h *Handler) CreateFixture(w http.ResponseWriter, r *http.Request) {
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

	var req CreateFixtureRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Inwalid request")
		return
	}

	if strings.ToLower(req.Type) != "fixture" {
		http.Error(w, "invalid request body: type must be Fixture", http.StatusBadRequest)
		return
	}

	if req.FixtureNo == "" {
		response.Error(w, http.StatusBadRequest, "Fixture no required")
		return
	}

	id, err := h.Service.Create(ctx, &req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}
