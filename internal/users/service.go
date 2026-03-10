package users

import (
	"context"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/utils"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateUser(ctx context.Context, req UserCreateRequest) (*UserResponse, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	// Call store
	user, err := ser.Store.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// convert DB model -> Response DTO

	res := &UserResponse{
		ID:         user.ID,
		TenantID:   user.TenantID,
		RoleID:     user.RoleID,
		EmployeeID: user.EmployeeID,
		UserName:   user.UserName,
		Phone:      utils.SafeString(user.Phone),
		Email:      utils.SafeString(user.Email),
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedBy:  utils.SafeInt(user.CreatedBy),
		UpdatedBy:  utils.SafeInt(user.UpdatedBy),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	return res, nil

}

func (ser *Service) LoginUser(ctx context.Context, req LoginRequest) (*UserResponse, error) {

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		return nil, err
	}

	TokenPayload, hashedPassword, err := ser.Store.GetPasswordHashbyEmplopeeID(ctx, req.EmployeeID)

	if err != nil {
		return nil, err
	}

	// compare password
	if err := utils.CompareHash(hashedPassword, req.Password); err != nil {
		return nil, err
	}

	authToken, err := auth.GenerateToken(TokenPayload.UserID, TokenPayload.TenantID)
	if err != nil {
		return nil, err
	}

}
