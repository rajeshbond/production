package users

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type Module struct {
	handler   *Handler
	service   *Service
	store     *Store
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth, roleProvider RoleProvide) *Module {
	store := NewStore(db)
	service := NewService(store, roleProvider)
	handler := NewHandler(service, tokenAuth)

	return &Module{
		store:     store,
		service:   service,
		handler:   handler,
		tokenAuth: tokenAuth,
	}
}
