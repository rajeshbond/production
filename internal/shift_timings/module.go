package shifttiming

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/database"
)

type Module struct {
	Handler   *Handler
	Service   *Service
	Store     *Store
	tokenAuth *jwtauth.JWTAuth
}

func NewModule(db *database.DB, tokenAuth *jwtauth.JWTAuth, roleProvider RoleProvider) *Module {

	// 🔹 Initialize Store (sql + pgx)
	store := NewStore(db)

	// 🔹 Initialize Service
	service := NewService(store, roleProvider)

	// 🔹 Initialize Handler
	handler := NewHandler(service, tokenAuth)

	return &Module{
		Handler:   handler,
		Service:   service,
		Store:     store,
		tokenAuth: tokenAuth,
	}
}
