package application

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/rajesh_bond/production/cmd/service"
	internalsetup "github.com/rajesh_bond/production/internal/internal_setup"
	tenant "github.com/rajesh_bond/production/internal/tenant"
	userrole "github.com/rajesh_bond/production/internal/user_role"
	users "github.com/rajesh_bond/production/internal/users"

	_ "github.com/rajesh_bond/production/docs"
)

func NewRouter(app *App) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	tokenAuth := service.GetTokenAuth()

	// Static files
	staticPath, _ := filepath.Abs("./web/static")
	fileServer := http.FileServer(http.Dir(staticPath))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Home
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath, _ := filepath.Abs("./web/index.html")
		http.ServeFile(w, r, indexPath)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health Ok"))
	})

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// ---------------- MODULES ----------------

	// Internal setup
	internalModule := internalsetup.NewModule(app.DB)
	r.Mount("/internal", internalModule.Router())

	// User Role
	userRoleModule := userrole.NewModule(app.DB)
	r.Mount("/user-role", userRoleModule.Router())

	// Tenant
	tenantModule := tenant.NewModule(app.DB)
	r.Mount("/tenant", tenantModule.Router())

	// Users
	usersModule := users.NewModule(app.DB, tokenAuth, userRoleModule.Service, tenantModule.Service)
	r.Mount("/users", usersModule.Router())

	return r
}
