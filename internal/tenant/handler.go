package tenant

// Index - Handler
//////////////////////////////////////

// 1. Create Teanant

// 2. Tenant Verification

//////////////////////////////////////

//////////////////////////////////////
// Code Starts Here
//////////////////////////////////////

// Imports

import (
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
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
	
	tenantCode := chi.URLParam(r,"tenant_code")

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
	
	okDeleted, err := h.service.DeleteTenant(ctx,tenantCode, claims.UserID)

	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	response.JSON(w, http.StatusOK, okDeleted)
}
