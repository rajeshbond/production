package operations

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
	"github.com/rajesh_bond/production/internal/common/utils"
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

func (h *Handler) CreateOperations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {

		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}

	if !auth.IsTenatAdminRole(claims.Role) {
		response.Error(w, http.StatusForbidden, auth.ErrUnauthorized.Error())
		return
	}

	var req BulkCreateOperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}
	if err := utils.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.Service.CreateOperations(ctx, req, claims)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, result)

}
