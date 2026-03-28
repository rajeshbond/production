package tenantshifts

import (
	"encoding/json"
	"errors"
	"log"
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

func (h *Handler) CreataTenantShiftResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Extract jwt claims from context

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusUnauthorized, "No JWT claims found in the context")
		return
	}

	var req CreateTenantShiftRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.Service.CreateTenantShift(ctx, claims, &req)
	if err != nil {

		switch {
		case errors.Is(err, ErrTenantShiftAlreadyExists):
			response.Error(w, http.StatusConflict, err.Error())

		case errors.Is(err, ErrInvalidRequest):
			response.Error(w, http.StatusBadRequest, err.Error())

		default:
			// log actual error
			log.Println("CreateTenantShift error:", err)

			response.Error(w, http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}
