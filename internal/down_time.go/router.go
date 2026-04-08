package downtime

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/downtime", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Test Ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(m.tokenAuth))
		r.Use(auth.Authenticator(m.tokenAuth))
		r.Use(auth.UserContextInjector)
		r.Post("/createdowntime", m.Handler.CreateDownTime)

	})

	return r
}
