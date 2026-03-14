package userrole

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/cmd/service"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	tokenAuth := service.GetTokenAuth()

	// Public routes

	// r.Post("/create", m.handler.Create)

	r.Group(func(r chi.Router) {

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Health Ok"))
		})

		r.Get("/test1", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test1 Ok"))
		})

		// r.Post("/createrole", m.handler.Create)

	})

	// Protected routes

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(tokenAuth))
		r.Use(auth.Authenticator(tokenAuth))
		r.Use(auth.UserContextInjector)
		// Create User role
		r.Post("/createrole", m.Handler.Create)
		r.Post("/test", m.Handler.TestRole1)
	})

	return r

}
