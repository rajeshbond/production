package tenant

import (
	"database/sql"
)

type Module struct {
	Handler *Handler
	Service *Service
	Store   *Store
}

func NewModule(db *sql.DB) *Module {
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service)

	return &Module{
		Store:   store,
		Service: service,
		Handler: handler,
	}
}
