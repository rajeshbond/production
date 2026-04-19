package shifttiming

import "context"

type RoleProvider interface {
	GetRoleNameByID(ctx context.Context, roleID int64) (string, error)
	GetRoleIDByName(ctx context.Context, roleName string) (int64, error)
}
