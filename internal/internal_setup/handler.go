package internalsetup

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupSuperAdmin(w http.ResponseWriter, r *http.Request) {

	// Allow only POST method
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var dto SetupSuperAdminDTO

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Println("[SETUP ERROR] invalid request body:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Normalize role name
	dto.Role.UserRole = strings.ToLower(dto.Role.UserRole)

	log.Println("[DEV SETUP REQUEST]", dto)

	// Call service layer
	result, err := h.service.SetupSuperAdmin(r.Context(), &dto)
	if err != nil {
		log.Println("[SETUP ERROR]", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("[SETUP SUCCESS]", result)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println("[RESPONSE ERROR]", err)
	}
}
