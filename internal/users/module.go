package users

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type Module struct {
	Handler   *Handler
	service   *Service
	store     *Store
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth, roleProvider RoleProvider, trnantProvide TenantProvider) *Module {
	store := NewStore(db)
	service := NewService(store, roleProvider, trnantProvide)
	handler := NewHandler(service, tokenAuth)

	return &Module{
		store:     store,
		service:   service,
		Handler:   handler,
		tokenAuth: tokenAuth,
	}
}
