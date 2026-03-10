package internalsetup

type RoleDTO struct {
	UserRole string `json:"user_role"`
}

type TenantDTO struct {
	TenantName string `json:"tenant_name"`
	TenantCode string `json:"tenant_code"`
	Address    string `json:"address"`
}

type UserDTO struct {
	EmployeeID string `json:"employee_id"`
	UserName   string `json:"user_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type SetupSuperAdminDTO struct {
	Role   RoleDTO   `json:"role"`
	Tenant TenantDTO `json:"tenant"`
	User   UserDTO   `json:"user"`
}
