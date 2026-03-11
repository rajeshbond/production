package internalsetup

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func devOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if os.Getenv("APP_ENV") != "dev" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Module) Router() chi.Router {

	r := chi.NewRouter()

	// Allow only in dev environment
	r.Use(devOnly)

	h := NewHandler(m.Service)

	r.Post("/setup-superadmin", h.SetupSuperAdmin)

	return r
}
