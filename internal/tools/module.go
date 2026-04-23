package tools

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

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth, resourcePointer ResourceProvider) *Module {
	store := NewStore(db)
	service := NewService(store, resourcePointer)
	handler := NewHandler(service)

	return &Module{
		tokenAuth: tokenAuth,
		Store:     store,
		Service:   service,
		Handler:   handler,
	}
}
