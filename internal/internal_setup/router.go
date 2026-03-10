package internalsetup

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func devOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("APP_ENV") != "dev" {
			http.Error(w, "Method not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(devOnly)
	h := NewHandler(m)
	r.Post("/setup-superadmin", h.SetupSuperAdmin)

	return r
}
