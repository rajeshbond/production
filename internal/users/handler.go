package users

import (
	"encoding/json"
	"net/http"

	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
