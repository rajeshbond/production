package tenant

// Index - Handler
//////////////////////////////////////

// 1. Create Teanant

// 2. Tenant Verification

// 3. Delete Tenant

//////////////////////////////////////

//////////////////////////////////////
// Code Starts Here
//////////////////////////////////////

// Imports

import (
	"encoding/json"
	"errors"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

// structs

type Handler struct {
	service *Service
}

// Struct connstructor
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Validator
var validate = validator.New()

// 1. Create Teanant
func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var req CreateTenantDTO

	// Get JWT claims from context
	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// user id from token
	role := claims.Role

	if !auth.IsSuper(role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// Decode Json safely

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode request
	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.InvalidRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create tenant

	tenant, err := h.service.CreateTenant(ctx, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, tenant)

}

// 2. Tenant Verification
func (h *Handler) VerifyTenant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req IsVerfiedRequest

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// user id from token
	role := claims.Role

	if !auth.IsSuper(role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// Decode Json safely

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode request
	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.InvalidRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.TenantVerifcation(ctx, req.TenantCode)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
	}

	response.JSON(w, http.StatusOK, resp)
}

// 3. Delete Tenant

func (h *Handler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get tenant code from request

	tenantCode := chi.URLParam(r, "tenant_code")

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// user id from token
	role := claims.Role

	if !auth.IsSuper(role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// Decode Json safely

	okDeleted, err := h.service.DeleteTenant(ctx, tenantCode, claims.UserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.JSON(w, http.StatusOK, okDeleted)
}

// 4. Update Tenant

func (h *Handler) UpdateTenant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// user id from token
	role := claims.Role

	if !auth.IsSuper(role) {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// 1. Get tenant_code from URL param
	tenantCode := chi.URLParam(r, "tenant_code")
	if tenantCode == "" {
		response.Error(w, http.StatusBadRequest, TenantCodeRequired)
		return
	}

	// 2. Decode request body
	var req UpdateTenantDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Ensure at least one field is provided (PATCH validation)
	if req.TenantName == nil &&
		req.ContactPersonName == nil &&
		req.ContactPhone == nil &&
		req.ContactEmail == nil &&
		req.Address == nil &&
		req.IsActive == nil &&
		req.UpdatedBy == nil {

		response.Error(w, http.StatusBadRequest, "At least one field is required")
		return
	}

	// 4. Call service
	updated, err := h.service.UpdateTenant(ctx, tenantCode, req)
	if err != nil {

		switch {
		case errors.Is(err, ErrTenantCodeRequired):
			response.Error(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, ErrTenantCodeNotFount):
			response.Error(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrTenantNotUpdated):
			response.Error(w, http.StatusConflict, err.Error())

		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 5. Success response
	resp := map[string]interface{}{
		"success": updated,
		"message": "tenant updated successfully",
	}

	response.JSON(w, http.StatusOK, resp)
}
