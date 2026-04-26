package prodlog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateProductionLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	id, err := h.Service.CreateProdctionLog(ctx, &req, claims)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}

// Soft Delete

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	logidStr := chi.URLParam(r, "id")
	logid, _ := strconv.ParseInt(logidStr, 10, 64)

	if err := h.Service.SoftDelete(ctx, logid, claims); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "Deleted Sucessfully")

}
