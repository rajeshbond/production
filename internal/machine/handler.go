package machine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
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

func (h *Handler) CreateMachine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	fmt.Println(claims.Role)

	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, response.NotAuthorized)
		return
	}

	var req CreateMachineRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.Service.CreateMachine(ctx, req, claims)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}

	response.JSON(w, 201, res)

}

func (h *Handler) GetAllMachines(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	claims, _ := auth.GetUserClaimsFromContext(ctx)

	res, err := h.Service.Store.GetAllMachineByTenantID(ctx, claims.TenantID)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}

	response.JSON(w, 200, res)
}

func (h *Handler) GetMachine(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ctx := r.Context()
	claims, _ := auth.GetUserClaimsFromContext(ctx)

	res, err := h.Service.GetMachineBYID(ctx, claims.TenantID, id)
	if err != nil {
		response.Error(w, 404, err.Error())
		return
	}

	response.JSON(w, 200, res)
}

func (h *Handler) UpdateMachine(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	claims, _ := auth.GetUserClaimsFromContext(ctx)

	var req UpdateMachineRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.Service.UpdateMachine(ctx, &req, claims)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}

	response.JSON(w, 200, "updated")
}
