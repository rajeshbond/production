package users

/*
Handler Index
/////////////////////////////////
// 1. Create User
// 2. Login
// 3. Test Private route
// 4. Create Tenant User
// 5. Verify Tenant User
// 6. Delete Tenant User

*/

// Imports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

// structs
type Handler struct {
	Service   *Service
	tokenAuth *jwtauth.JWTAuth
}

// Struct connstructor
func NewHandler(service *Service, tokenAuth *jwtauth.JWTAuth) *Handler {
	return &Handler{
		Service:   service,
		tokenAuth: tokenAuth,
	}
}

// 1. Create Tenant Admin
func (h *Handler) CreateTenantAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	makeTenantAdmin := auth.IsSuper(claims.Role)

	// fmt.Println("makeTenantAdmin", makeTenantAdmin)

	if !makeTenantAdmin {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	var req UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body ", http.StatusBadRequest)
		return
	}

	user, err := h.Service.CreateTenantAdmin(ctx, claims, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, user)

}

// 2. Login
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var req LoginRequest

	// decode request body

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body ", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := h.Service.LoginUser(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, resp)

}

// 3. Test Private route
func (h *Handler) Test1(w http.ResponseWriter, r *http.Request) {

	// Extract JWT claims from context
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "No JWT claims found in context")
		return
	}

	resp := map[string]interface{}{
		"user_id":     claims.UserID,
		"employee_id": claims.EmployeeID,
		"tenant_id":   claims.TenantID,
		"role_id":     claims.RoleID,
		"username":    claims.Username,
		"message":     "Private route working",
	}
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	// jsonData, err := json.Marshal(resp)
	// if err != nil {
	// 	log.Println("JSON marshal error:", err)
	// 	return
	// }

	log.Println(string(jsonData))

	// log.Println(resp)
	response.JSON(w, http.StatusOK, resp)
	// response.JSON(w, http.StatusOK, string(jsonData))
}

// 4. Create Tenant User
func (h *Handler) CreateTenantUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	// Extrating claims (Only Here)
	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, "Permisson not allowed , only tenant Admin can perform operations...")
		return
	}
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// Parse Request
	var req UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.InvalidRequestBody)
		return
	}

	if claims.TenantID != req.TenantID {
		response.Error(w, http.StatusForbidden, "Not allowed for this tenant")
		return
	}

	// Call Service

	resp, err := h.Service.CreateTenantUser(ctx, claims, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, resp)

}

// 5. Verify Tenant User
func (h *Handler) VerifyTenantUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	// Extraxt claims

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		fmt.Println("Rajesh Bondgilwar")
		response.Error(w, http.StatusUnauthorized, response.NotAuthorized)
		return
	}

	// Pasrse request
	var req VerifyTenantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	// call service

	err := h.Service.VerifyTenantUser(ctx, claims, req.EmployeeID, req.TenantID)

	// Better error handleing (optional but recommended)
	if err != nil {

		// 🔥 Better error handling (optional but recommended)
		switch err.Error() {
		case "user not found":
			response.JSON(w, http.StatusNotFound, err.Error())
		case "user already verified":
			response.JSON(w, http.StatusConflict, err.Error())
		default:
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, "user verified sucessfully ")

}

// 6. Delete Tenant User
func (h *Handler) DeleteTenantUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get params from URL

	employeeID := chi.URLParam(r, "employee_id")
	tenantIDStr := chi.URLParam(r, "tenant_id")

	// Validate tenant_id

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tenant id",
		})
		return
	}

	// Get claims from conetxt

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid tenant_id",
		})
		return
	}

	// call service

	err = h.Service.DeleteTenantUser(ctx, claims, employeeID, tenantID)

	if err != nil {

		switch err.Error() {
		case "employee_id is required", "tenant_id is required":
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		case "user not found":
			response.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})

		case "user already deleted":
			response.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": err.Error(),
			})

		default:
			response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})

		}

		return

	}

	// ✅ Success
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})

}
