package internalsetup

import (
	"database/sql"

	"github.com/rajesh_bond/production/internal/tenant"
	userrole "github.com/rajesh_bond/production/internal/user_role"
	"github.com/rajesh_bond/production/internal/users"
)

type Module struct {
	Service *Service
}

func NewModule(db *sql.DB) *Module {

	// Initialize stores
	roleStore := userrole.NewStore(db)
	tenantStore := tenant.NewStore(db)
	userStore := users.NewStore(db)

	// Initialize services
	roleService := userrole.NewService(roleStore)
	tenantService := tenant.NewService(tenantStore)
	userService := users.NewService(userStore, users.RoleProvider(roleService), tenantStore)

	// Initialize setup service
	setupService := NewService(db, tenantService, roleService, userService)

	return &Module{
		Service: setupService,
	}
}
