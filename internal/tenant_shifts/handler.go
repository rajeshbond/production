package tenantshifts

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

func (h *Handler) CreateTenanthifts(w http.ResponseWriter, r *http.Request) {
	// Authentication
	ctx := r.Context()

	// Extract jwt claims from context

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
	}

	// Decode request

	var req CreateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.InvalidRequest)
	}

	if err := auth.ValidateTenantAccesswithTenantCode(claims.Role, claims.TenantID, req.TenantID); err != nil {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
	}
	// Call Sevice layer

	resp, err := h.Service.CreateTenantShifts(ctx, req, claims.UserID)

	if err != nil {
		// You can improve error handling later (custom errors)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sucess response

	response.JSON(w, http.StatusCreated, resp)

}

// func (h *Handler) CreataTenantShiftResponse(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	// Extract jwt claims from context

// 	claims, ok := auth.GetUserClaimsFromContext(ctx)

// 	if !ok {
// 		response.Error(w, http.StatusUnauthorized, "No JWT claims found in the context")
// 		return
// 	}

// 	var req CreateTenantShiftRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)

// 	resp, err := h.Service.CreateTenantShift(ctx, claims, &req)
// 	if err != nil {

// 		switch {
// 		case errors.Is(err, ErrTenantShiftAlreadyExists):
// 			response.Error(w, http.StatusConflict, err.Error())

// 		case errors.Is(err, ErrInvalidRequest):
// 			response.Error(w, http.StatusBadRequest, err.Error())

// 		default:
// 			// log actual error
// 			log.Println("CreateTenantShift error:", err)

// 			response.Error(w, http.StatusInternalServerError, "Something went wrong")
// 		}
// 		return
// 	}

// 	response.JSON(w, http.StatusCreated, resp)
// }
