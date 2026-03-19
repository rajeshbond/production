package userrole

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var dto CreateUserRoleDTO

	// Decode JSON safely
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate required field
	if dto.UserRole == "" {
		response.Error(w, http.StatusBadRequest, "user_role is required")
		return
	}

	ctx := r.Context()

	// Get JWT claims from middleware
	claims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := strings.ToLower(claims.Role)
	dto.UserRole = strings.ToLower(dto.UserRole)

	fmt.Println("Role:-", role)

	if !auth.IsSuper(role) {
		response.Error(w, http.StatusUnauthorized, "Yoar and not authorized")
	}

	createdRole, err := h.service.Create(ctx, dto)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Assign created_by
	// dto.CreatedBy = &userID

	// Call service
	// role, err := h.service.Create(ctx,dto)
	// if err != nil {
	// 	response.Error(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	response.JSON(w, http.StatusCreated, createdRole)
}

func (h *Handler) TestRole1(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get JWT claims from middleware
	claims, _ := auth.GetUserClaimsFromContext(ctx)
	fmt.Println(claims)
	response.JSON(w, http.StatusOK, claims)

}

// CreateUserRole godoc
//	@Summary		Create User Role
//	@Description	Create a new user role
//	@Tags			UserRole
//	@Accept			json
//	@Produce		json
//	@Param			user_role	body		CreateUserRoleDTO	true	"User Role"
//	@Success		201			{object}	UserRoleResponseDTO
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Router			/user-role/createrole [post]
//	@Security		ApiKeyAuth

// func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

// 	var dto CreateUserRoleDTO

// 	// ✅ Decode JSON safely
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	if err := decoder.Decode(&dto); err != nil {
// 		response.Error(w, http.StatusBadRequest, "Invalid JSON body")
// 		return
// 	}

// 	// ✅ Validate required field
// 	if dto.UserRole == "" {
// 		response.Error(w, http.StatusBadRequest, "user_role is required")
// 		return
// 	}

// 	// ✅ Get user ID from middleware context
// 	user := auth.GetUserIDFromContext(r.Context())

// 	if userID != "" {

// 		id, err := strconv.ParseInt(userID, 10, 64)
// 		fmt.Println("Handdler", id)
// 		if err != nil {
// 			response.Error(w, http.StatusBadRequest, "Invalid user ID in token")
// 			return
// 		}
// 		dto.CreatedBy = &id
// 		fmt.Printf("DTO Handler :%+v", dto)
// 	} else {

// 		id := int64(1)

// 		dto.CreatedBy = &id

// 		fmt.Printf("DTO Handler else :%+v", dto)

// 	}

// 	// ✅ Call service
// 	role, err := h.service.Create(r.Context(), dto)
// 	if err != nil {
// 		response.Error(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	response.JSON(w, http.StatusCreated, role)
// }

// func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

// 	var dto CreateUserRoleDTO

// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		response.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	// get the user id from context (middleware)

// 	userID := auth.GetUserIDFromContext(r.Context())

// 	// if user exist -> set createdBy
// 	if userID != "" {
// 		id, err := strconv.ParseInt(userID, 10, 64)

// 		if err != nil {
// 			log.Println("test line")
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		dto.CreatedBy = &id
// 	}

// 	role, err := h.service.Creare(r.Context(), dto)
// 	if err != nil {
// 		response.Error(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	response.JSON(w, http.StatusCreated, role)

// }
