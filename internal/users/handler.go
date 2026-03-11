package users

import (
	"encoding/json"
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

	// Set content type and status manually
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusInternalServerError, "No user ID found in context")
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Server working via Mount"))
	response := map[string]interface{}{
		"user_id": userID,
		"message": "Private route working",
	}
	_ = json.NewEncoder(w).Encode(response)

}
