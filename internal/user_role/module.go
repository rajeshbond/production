package userrole

import (
	"database/sql"
)

type Module struct {
	Handler *Handler // Capitalized
	Service *Service // Capitalized
	Store   *Store   // Capitalized
}

func NewModule(db *sql.DB) *Module {
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service)

	return &Module{
		Handler: handler, // Use uppercase
		Service: service, // Use uppercase
		Store:   store,   // Use uppercase
	}
}
