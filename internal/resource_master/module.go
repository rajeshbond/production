package resourcemaster

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type Module struct {
	Store     *Store
	Servie    *Service
	Handler   *Handler
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *sql.DB, tokenAuth *jwtauth.JWTAuth) *Module {
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service, tokenAuth)

	return &Module{
		tokenAuth: tokenAuth,
		Store:     store,
		Servie:    service,
		Handler:   handler,
	}
}
