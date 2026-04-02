package users

import "context"

type RoleProvider interface {
	GetRoleNameByID(ctx context.Context, roleID int64) (string, error)
}
