package service

type TokenPayload struct {
	TenantID int64  `json:"tenant_id"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
}
