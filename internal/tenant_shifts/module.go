package tenantshifts

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

func NewModule(db *sql.DB, tokenauth *jwtauth.JWTAuth) *Module {

	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service, tokenauth)

	return &Module{
		tokenAuth: tokenauth,
		store:     store,
		service:   service,
		handler:   handler,
	}
}
