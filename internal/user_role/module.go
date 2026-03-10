package userrole

import (
	"database/sql"
)

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
		handler: handler,
		service: service,
		store:   store,
	}
}

// func (m *Module) Router() chi.Router {
// 	r := chi.NewRouter()
// 	tokenAuth := service.GetTokenAuth()

// 	// Public routes

// 	// r.Post("/create", m.handler.Create)

// 	r.Group(func(r chi.Router) {

// 		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte("Health Ok"))
// 		})

// 		r.Get("/test1", func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte("test1 Ok"))
// 		})

// 	})

// 	// r.Post("/create", m.handler.Create)
// 	// r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
// 	// 	w.Write([]byte("Health Ok"))
// 	// })

// 	// Protected routes

// 	r.Group(func(r chi.Router) {
// 		r.Use(auth.Verifier(tokenAuth))
// 		r.Use(auth.Authenticator(tokenAuth))
// 		r.Use(auth.UserContextInjector)
// 	})

// 	return r

// }
