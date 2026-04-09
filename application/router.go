package application

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/rajesh_bond/production/cmd/service"
	"github.com/rajesh_bond/production/internal/defect"
	downtime "github.com/rajesh_bond/production/internal/down_time.go"
	internalsetup "github.com/rajesh_bond/production/internal/internal_setup"
	operationdefectmap "github.com/rajesh_bond/production/internal/operation_defect_map"
	operationdowntimemap "github.com/rajesh_bond/production/internal/operation_downtime_map"
	"github.com/rajesh_bond/production/internal/operations"
	shifttiming "github.com/rajesh_bond/production/internal/shift_timings"
	tenant "github.com/rajesh_bond/production/internal/tenant"
	tenantshifts "github.com/rajesh_bond/production/internal/tenant_shifts"
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
	internalModule := internalsetup.NewModule(app.DB.SQLDB)
	r.Mount("/internal", internalModule.Router())

	// User Role
	userRoleModule := userrole.NewModule(app.DB.SQLDB)
	r.Mount("/user-role", userRoleModule.Router())

	// Tenant
	tenantModule := tenant.NewModule(app.DB.SQLDB)
	r.Mount("/tenant", tenantModule.Router())

	// Tenant shift
	tenantShiftMoulde := tenantshifts.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/tenantshift", tenantShiftMoulde.Router())

	// Shift timings

	shiftTimingsModule := shifttiming.NewModule(app.DB, tokenAuth, tenantModule.Store)
	r.Mount("/shifttiming", shiftTimingsModule.Router())

	// Defect

	defectModule := defect.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/defect", defectModule.Router())

	// Downtime

	downtimeNodule := downtime.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/downtime", downtimeNodule.Router())

	// Operation

	operationModule := operations.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/operations", operationModule.Router())

	// Operation Defect Map

	operationDefectMap := operationdefectmap.NewModule(app.DB.SQLDB, tokenAuth, defectModule.Store, operationModule.Store)
	r.Mount("/opdefmap", operationDefectMap.Router())

	// Operation Downtime Map

	operationDowntimeMap := operationdowntimemap.NewModule(app.DB.SQLDB, tokenAuth, downtimeNodule.Store, operationModule.Store)
	r.Mount("/mapopdt", operationDowntimeMap.Router())

	// Users
	usersModule := users.NewModule(app.DB.SQLDB, tokenAuth, userRoleModule.Service, tenantModule.Service)
	r.Mount("/users", usersModule.Router())

	return r
}
