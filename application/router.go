package application

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	internalsetup "github.com/rajesh_bond/production/internal/internal_setup"
	tenant "github.com/rajesh_bond/production/internal/tenant"
	userrole "github.com/rajesh_bond/production/internal/user_role"
	"github.com/rajesh_bond/production/internal/users"

	_ "github.com/rajesh_bond/production/docs"
)

func NewRouter(app *App) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// static files
	staticPath, _ := filepath.Abs("./web/static")
	fileServer := http.FileServer(http.Dir(staticPath))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Home
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath, _ := filepath.Abs("./web/index.html")
		http.ServeFile(w, r, indexPath)
	})

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health Ok"))
	})

	// ✅ Swagger Route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Modules ------------->

	// Internal Modules
	internalModule := internalsetup.NewModule(app.DB)
	r.Mount("/internal", internalModule.Router())
	// User role
	userRole := userrole.NewModule(app.DB)
	r.Mount("/user-role", userRole.Router())
	// tenant
	tenantModule := tenant.NewModule(app.DB)
	r.Mount("/tenant", tenantModule.Router())
	// Users

	users := users.NewModule(app.DB)
	r.Mount("/users", users.Router())

	return r
}
