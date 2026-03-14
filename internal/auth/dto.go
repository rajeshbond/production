package auth

type UserClaims struct {
	EmployeeID string `json:"employee_id"`
	Iat        int64  `json:"iat"`
	RoleID     int64  `json:"role_id"`
	Role       string `json:"role"`
	TenantID   int64  `json:"tenant_id"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
}

// type UserClaims struct {
// 	EmployeeID string `json:"employee_id"`
// 	Iat        int64  `json:"iat"`
// 	RoleID     int64  `json:"role_id"`
// 	TenantID   int64  `json:"tenant_id"`
// 	UserID     int64  `json:"user_id"`
// 	Username   string `json:"username"`
// }
