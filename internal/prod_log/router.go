package prodlog

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/test1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ProductionLog Test Ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(m.tokenAuth))
		r.Use(auth.Authenticator(m.tokenAuth))
		r.Use(auth.UserContextInjector)
		r.Post("/createpol", m.Handler.CreateProductionLog)
		r.Delete("/pol/{id}", m.Handler.Delete)

	})

	return r
}
