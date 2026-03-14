package users

import "context"

type RoleProvide interface {
	GetRoleNameByID(ctx context.Context, roleID int64) (string, error)
}
