package auth

import (
	"context"

	"github.com/Template7/backend/internal/t7Error"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
)

func (s *service) GetUserRole(ctx context.Context, userId string) (authV1.Role, error) {
	roles, err := s.core.GetRolesForUser(userId)
	if err != nil {
		return 0, t7Error.DbOperationFail.WithDetail(err.Error())
	}
	if len(roles) == 0 {
		return 0, t7Error.UserHasNoRole
	}
	if len(roles) > 1 {
		s.log.WithContext(ctx).With("userId", userId).With("roles", roles).Warn("user has multiple roles, using the first")
	}

	role, ok := authV1.Role_value[roles[0]]
	if !ok {
		return 0, t7Error.UserHasNoRole.WithDetail("unknown role " + roles[0])
	}
	return authV1.Role(role), nil
}

func (s *service) AddUserRole(_ context.Context, userId string, role authV1.Role) error {
	if _, err := s.core.AddRoleForUser(userId, role.String()); err != nil {
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}
	return nil
}

func (s *service) DeleteUserRoles(_ context.Context, userId string) error {
	if _, err := s.core.DeleteRolesForUser(userId); err != nil {
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}
	return nil
}
