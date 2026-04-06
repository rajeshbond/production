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

func (h *Handler) BulkCreateShift(w http.ResponseWriter, r *http.Request) {
	var req BulkCreateShiftRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusUnauthorized, auth.UnAuthorised)
		return
	}

	result, err := h.Service.BulkCreateShift(ctx, req, claims)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// response.JSON(w, http.StatusCreated, result)
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"result": result,
		"remark": "created",
	})

}
