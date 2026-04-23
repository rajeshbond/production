package mould

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type Module struct {
	Store     *Store
	Service   *Service
	Handler   *Handler
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth, resourceProvide ResourceProvider) *Module {
	store := NewStore(db)
	service := NewService(store, resourceProvide)
	handler := NewHandler(service)

	return &Module{
		Store:     store,
		Service:   service,
		Handler:   handler,
		tokenAuth: tokenAuth,
	}
}
