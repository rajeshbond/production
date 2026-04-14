package shiftslot

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajesh_bond/production/internal/auth"
)

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/shiftslot", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Test Ok"))
	})
	// r.Get("/shifttime", m.Handler.GenerateAllShiftSlots)

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(m.tokenAuth))
		r.Use(auth.Authenticator(m.tokenAuth))
		r.Use(auth.UserContextInjector)
		r.Get("/shifttime1", m.Handler.GenerateAllShiftSlots)

	})

	return r
}
