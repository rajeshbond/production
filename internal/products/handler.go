package products

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}
	if claims.Role != "tenantadmin" {
		response.Error(w, http.StatusForbidden, auth.ErrUnauthorized.Error())
		return
	}

	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.Service.CreateProduct(ctx, req.ProductName, req.ProductNo, claims)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, id)

}

// Get Product By ID

func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusForbidden, auth.ErrForbidden.Error())
		return
	}

	if claims.Role != "tenantadmin" {
		if claims.Role != "tenantadmin" {
			response.Error(w, http.StatusForbidden, auth.ErrUnauthorized.Error())
			return
		}
	}

	idStr := r.URL.Query().Get("id")
	productID, _ := strconv.ParseInt(idStr, 10, 64)

	name, number, err := h.Service.GetProductID(ctx, productID, claims)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           productID,
		"product_name": name,
		"product_no":   number,
	})
}

func (h *Handler) SearchProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get tenant_if from path

	// // tenantIDStr := chi.URLParam(r, "tenant_id")
	// // if tenantIDStr == "" {
	// // 	response.Error(w, http.StatusBadRequest, "tenant_id is required")
	// // 	return
	// // }

	// // tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)

	// if err != nil {
	// 	response.Error(w, http.StatusBadRequest, "invalid tenant_id")
	// 	return
	// }

	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusBadRequest, auth.UnAuthorised)
		return
	}

	Products, err := h.Service.GetAllProductsByTenant(ctx, claims)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, Products)

}

func (h *Handler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.URL.Query().Get("q")

	claims, ok := auth.GetUserClaimsFromContext(ctx)

	if !ok {
		response.Error(w, http.StatusBadRequest, auth.UnAuthorised)
		return
	}

	products, err := h.Service.SerachProducts(ctx, search, claims)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, products)

}
