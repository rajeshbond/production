package operationdowntimemap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/mapdowntime", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Downtime Map test ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(m.tokenAuth))
		r.Use(auth.Authenticator(m.tokenAuth))
		r.Use(auth.UserContextInjector)
		r.Post("/mapDowntimeCreate", m.Handler.CreateOperationWithDowntime)

	})

	return r

}
