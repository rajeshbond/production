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
	"github.com/rajesh_bond/production/internal/fixture"
	internalsetup "github.com/rajesh_bond/production/internal/internal_setup"
	"github.com/rajesh_bond/production/internal/machine"
	"github.com/rajesh_bond/production/internal/mold"
	"github.com/rajesh_bond/production/internal/mould"
	operationdefectmap "github.com/rajesh_bond/production/internal/operation_defect_map"
	operationdowntimemap "github.com/rajesh_bond/production/internal/operation_downtime_map"
	operationmode "github.com/rajesh_bond/production/internal/operation_mode"
	"github.com/rajesh_bond/production/internal/operations"
	productoperationsequence "github.com/rajesh_bond/production/internal/product_operation_sequence"
	productionlog "github.com/rajesh_bond/production/internal/production_log"
	productiontarget "github.com/rajesh_bond/production/internal/production_target"
	"github.com/rajesh_bond/production/internal/products"
	"github.com/rajesh_bond/production/internal/resource"
	resourcemaster "github.com/rajesh_bond/production/internal/resource_master"
	resourcetype "github.com/rajesh_bond/production/internal/resource_type"
	"github.com/rajesh_bond/production/internal/shift"
	shiftslot "github.com/rajesh_bond/production/internal/shift_slot"
	shifttiming "github.com/rajesh_bond/production/internal/shift_timings"
	tenant "github.com/rajesh_bond/production/internal/tenant"
	tenantshifts "github.com/rajesh_bond/production/internal/tenant_shifts"
	"github.com/rajesh_bond/production/internal/tools"
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

	// Shift

	shiftModule := shift.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/shift", shiftModule.Router())

	// shift_slots

	shiftSlot := shiftslot.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/shiftslot", shiftSlot.Router())

	// Defect

	defectModule := defect.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/defect", defectModule.Router())

	// Downtime

	// Machine

	machineModule := machine.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/machine", machineModule.Router())

	downtimeNodule := downtime.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/downtime", downtimeNodule.Router())

	// Operation

	operationModule := operations.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/operations", operationModule.Router())

	// Operation mode

	operationModeModule := operationmode.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/opmode", operationModeModule.Router())

	// Operation Defect Map

	operationDefectMap := operationdefectmap.NewModule(app.DB.SQLDB, tokenAuth, defectModule.Store, operationModule.Store)
	r.Mount("/opdefmap", operationDefectMap.Router())

	// Product

	productModule := products.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/product", productModule.Router())

	operationDowntimeMap := operationdowntimemap.NewModule(app.DB.SQLDB, tokenAuth, downtimeNodule.Store, operationModule.Store)
	r.Mount("/mapopdt", operationDowntimeMap.Router())

	// Product Operation Sequence

	productOperationSequence := productoperationsequence.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/pos", productOperationSequence.Router())

	// Users
	usersModule := users.NewModule(app.DB.SQLDB, tokenAuth, userRoleModule.Service, tenantModule.Service)
	r.Mount("/users", usersModule.Router())

	// Resource Type

	resourcetypeModule := resourcetype.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/resty", resourcetypeModule.Router())

	// Resource Master

	resourceMasterModule := resourcemaster.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/resmaster", resourceMasterModule.Router())

	// Mould

	moldModule := mold.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/mold", moldModule.Router())

	// New

	mouldModule := mould.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/mould", mouldModule.Router())

	// Fixture

	fixtureModule := fixture.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/fixture", fixtureModule.Router())

	// Tools

	toolsModule := tools.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/tools", toolsModule.Router())

	// Resource

	resourceModule := resource.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/resource", resourceModule.Router())

	// Production Target

	productionTarget := productiontarget.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/prodtrmap", productionTarget.Router())

	// Production Log
	productionLogModule := productionlog.NewModule(app.DB.SQLDB, tokenAuth)
	r.Mount("/plog", productionLogModule.Router())

	return r
}
