package products

import (
	"context"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateProduct(ctx context.Context, productName string, productNumber string, claims *auth.UserClaims) (int64, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	id, err := ser.Store.CreateProduct(ctx, tx, claims.TenantID, productName, productNumber, claims.UserID)
	if err != nil {
		return 0, nil
	}

	if err := tx.Commit(); err != nil {
		return 0, nil
	}

	return id, nil

}

func (ser *Service) GetProductID(ctx context.Context, productID int64, claims *auth.UserClaims) (string, string, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	productName, productNumber, err := ser.Store.GetProductByID(ctx, tx, claims.TenantID, productID)
	if err != nil {
		return "", "", err
	}

	return productName, productNumber, nil

}

// Search Product by name or number

func (ser *Service) GetProductByNoOrName(ctx context.Context, value string, claims *auth.UserClaims) (ProductResponse, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return ProductResponse{}, nil
	}
	defer tx.Rollback()

	product, err := ser.Store.GetProductByNoOrName(ctx, tx, claims.TenantID, value)
	if err != nil {
		return ProductResponse{}, nil
	}

	return product, nil

}

func (ser *Service) SerachProducts(ctx context.Context, search string, claims *auth.UserClaims) ([]ProductResponse, error) {
	return ser.Store.SearchProducts(ctx, claims.TenantID, search)
}

func (ser *Service) UpdateProduct(ctx context.Context, req *UpdateProductRequest, claims *auth.UserClaims) error {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = ser.Store.UpdateProduct(ctx, tx, claims.TenantID, claims.UserID, req)
	if err != nil {
		return nil
	}

	return tx.Commit()
}

func (ser *Service) DeleteProduct(ctx context.Context, productID int64, claims *auth.UserClaims) error {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = ser.Store.DeleteProduct(ctx, tx, claims.TenantID, productID, claims.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (ser *Service) GetAllProductsByTenant(ctx context.Context, claims *auth.UserClaims) ([]Product, error) {
	produts, err := ser.Store.GetAllProductsByTenant(ctx, claims.TenantID)

	if err != nil {
		return nil, err
	}

	return produts, nil
}
