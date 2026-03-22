package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {

	r := chi.NewRouter()
	r.Get("/user-test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Test Ok"))
	})

	// Public routes
	r.Post("/createuser", m.handler.CreateUser)
	r.Post("/loginuser", m.handler.LoginUser)

	// Private routes

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(m.tokenAuth))
		r.Use(auth.Authenticator(m.tokenAuth))
		r.Use(auth.UserContextInjector)

		r.Get("/test", m.handler.Test1)
		r.Post("/ctenatuser", m.handler.CreateTenantUser)

	})
	return r
}
