package tenant

import (
	"database/sql"
)

type Module struct {
	Handler *Handler
	Service *Service
	store   *Store
}

func NewModule(db *sql.DB) *Module {
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service)

	return &Module{
		store:   store,
		Service: service,
		Handler: handler,
	}
}
