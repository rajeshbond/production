package shifttiming

import (
	"encoding/json"
	"net/http"

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
		Service:   service,
		tokenAuth: tokenAuth,
	}
}

func (h *Handler) CreateShiftTimings(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	var req CreateShiftTimingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusForbidden, response.InvalidRequestBody)
		return
	}

	if claims.Role != auth.RoleTenantAdmin {
		response.Error(w, http.StatusForbidden, response.OnlyTenantAllowed)
		return
	}

	resp, err := h.Service.CreateShiftTiming(ctx, req, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, resp)
}
