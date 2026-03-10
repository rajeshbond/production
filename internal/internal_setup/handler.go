package internalsetup

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	module *Module
}

func NewHandler(m *Module) *Handler {
	return &Handler{module: m}
}

func (h *Handler) SetupSuperAdmin(w http.ResponseWriter, r *http.Request) {
	var dto SetupSuperAdminDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// TODO: call
	// call service
	// creat role
	// tanent user

	// ✅ convert UserRole to lowercase
	dto.Role.UserRole = strings.ToLower(dto.Role.UserRole)

	log.Println("[DEV SETUP]", dto)
	w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Super admin created"))

}
