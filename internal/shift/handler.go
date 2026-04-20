package shift

import (
	"encoding/json"
	"net/http"
	"strings"

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

// ==========================================
// BULK CREATE (Multi Shift API)
// ==========================================
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// 🔐 Auth
	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// ✅ IMPORTANT: Accept ARRAY request
	var req BulkCreateShiftRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 🔥 strict validation

	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON body: "+err.Error())
		return
	}

	// Optional: empty check
	if len(req) == 0 {
		response.Error(w, http.StatusBadRequest, "Request body cannot be empty")
		return
	}

	userID := claims.UserID

	// 🚀 Call service
	err := h.Service.CreateBulk(ctx, userID, req)
	if err != nil {

		// DB overlap error
		if strings.Contains(err.Error(), "Shift overlap") {
			response.Error(w, http.StatusBadRequest, "Shift overlap detected")
			return
		}

		// Duplicate constraint
		if strings.Contains(err.Error(), "uix_shift_timing") {
			response.Error(w, http.StatusBadRequest, "Duplicate shift timing")
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, "Shifts created successfully")
}
