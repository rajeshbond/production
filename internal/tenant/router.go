package tenant

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health Ok"))
	})
	r.Post("/createtenant", m.Handler.CreateTenant)
	return r
}
