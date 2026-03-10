package users

import "database/sql"

type Module struct {
	handler *Handler
	service *Service
	store   *Store
}

func NewModule(db *sql.DB) *Module {
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service)

	return &Module{
		store:   store,
		service: service,
		handler: handler,
	}
}
