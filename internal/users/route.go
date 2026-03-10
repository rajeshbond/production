package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *Module) Router() chi.Router {

	r := chi.NewRouter()
	r.Get("/user-test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Test Ok"))
	})
	r.Post("/createuser", m.handler.CreateUser)
	return r
}
