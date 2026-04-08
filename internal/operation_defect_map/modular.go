package operationdefectmap

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type Module struct {
	Handler   *Handler
	Service   *Service
	Store     *Store
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth, defectProvider DefectProvider, operationProvider OperationProvider) *Module {
	store := NewStore(db)
	service := NewService(store, defectProvider, operationProvider)
	handler := NewHandler(service)

	return &Module{
		tokenAuth: tokenAuth,
		Handler:   handler,
		Service:   service,
		Store:     store,
	}
}
