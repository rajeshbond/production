package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	service   *Service
	tokenAuth *jwtauth.JWTAuth
}

func NewHandler(service *Service, tokenAuth *jwtauth.JWTAuth) *Handler {
	return &Handler{
		service:   service,
		tokenAuth: tokenAuth,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Rajesh Bondgilwar")
	var req UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body ", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := h.service.CreateUser(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response.JSON(w, http.StatusCreated, user)

}

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

	resp, err := h.service.LoginUser(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response.JSON(w, http.StatusOK, resp)

}

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
