package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

var validate = validator.New()

func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var req CreateTenantRequied

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get JWT claims from context
	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// user id from token
	userID := claims.UserID

	dto := CreateTenantDTO{
		TenantName: req.TenantName,
		TenantCode: req.TenantCode,
		Address:    req.Address,
		CreatedBy:  userID,
	}

	tenant, err := h.service.CreateTenant(ctx, dto)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, tenant)
}

// func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	var req CreateTenantRequied

// 	// ✅ Decode JSON properly
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.Error(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	if err := validate.Struct(req); err != nil {
// 		response.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	// ✅ Get user ID from middleware context
// 	userID := auth.GetUserIDFromContext(ctx)

// 	// if userID == "" {
// 	// 	response.Error(w, http.StatusUnauthorized, "unauthorized")
// 	// 	return
// 	// }
// 	var id int64
// 	if userID != "" {

// 		id, err := strconv.ParseInt(userID, 10, 64)
// 		fmt.Println("Handdler", id)
// 		fmt.Println("Inside if")
// 		if err != nil {
// 			response.Error(w, http.StatusBadRequest, "Invalid user ID in token")
// 			return
// 		}

// 		// fmt.Printf("DTO Handler :%+v", dto)
// 	} else {

// 		id = int64(1)
// 		// fmt.Println("=====inside else")

// 		// fmt.Printf("DTO Handler else :%+v", dto)

// 	}

// 	var dto = CreateTenantDTO{
// 		TenantName: req.TenantName,
// 		TenantCode: req.TenantCode,
// 		Address:    req.Address,
// 		CreatedBy:  id,
// 	}

// 	// fmt.Printf("DTO Handler: %+v\n", dto)

// 	// ✅ Call service correctly
// 	tenant, err := h.service.CreateTenant(ctx, dto)
// 	if err != nil {
// 		response.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	// // ✅ Success response
// 	response.JSON(w, http.StatusCreated, tenant)
// }
