package shiftslot

import (
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

func (h *Handler) GenerateAllShiftSlots(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}
	// var req GenerateSlotRequest

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }

	err := h.Service.GenerateAllShiftSlots(r.Context(), claims.TenantID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("All shift slots generated successfully"))
}
