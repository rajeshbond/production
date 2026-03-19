package auth

type Role string

const (
	RoleSuperAdmin  Role = "superadmin"
	RoleAdmin       Role = "admin"
	RoleTenantAdmin Role = "tenantadmin"
	RoleUser        Role = "user"
	RoleViewer      Role = "viewer"
)
