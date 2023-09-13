package auth

import (
	"context"
	"github.com/Template7/backend/internal/pkg/t7Error"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
)

func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
	log := s.log.WithContext(ctx)
	log.With("user", username).Debug("user login")

	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		log.WithError(err).Error("fail to get user")
		return "", t7Error.DbOperationFail.WithDetail(err.Error())
	}

	if !checkPasswordHash(password, user.Password) {
		log.Info("password incorrect")
		return "", t7Error.PasswordIncorrect
	}
	role := s.GetUserRole(ctx, username)
	if _, ok := v1.Role_name[int32(role)]; !ok {
		log.With("role", role).Warn("invalid user role")
		return "", t7Error.UserHasNoRole
	}

	token, err := s.genUserToken(ctx, username, s.GetUserRole(ctx, username))
	if err != nil {
		log.WithError(err).Error("fail to generate user token")
		return "", err
	}
	return token, nil
}

func (s *service) CheckPermission(ctx context.Context, sub, obj, act string) bool {
	log := s.log.WithContext(ctx).With("sub", sub).With("obj", obj).With("act", act)
	log.Debug("check permission")

	ok, _ := s.core.Enforce(sub, obj, act)
	if !ok {
		log.Info("no permission")
		return false
	}

	log.Debug("permission check ok")
	return true
}

func (s *service) GetUserRole(ctx context.Context, username string) v1.Role {
	log := s.log.WithContext(ctx).With("username", username)
	log.Debug("get user role")

	roles, err := s.core.GetRolesForUser(username)
	if err != nil {
		log.WithError(err).Error("unable to get user role")
		return -1
	}
	if len(roles) == 0 {
		log.Warn("user have no roles")
		return -1
	}
	if len(roles) > 1 {
		log.With("roles", roles).Warn("user has multiple roles")
	}

	log.With("role", roles).Debug("got user role")
	return v1.Role(v1.Role_value[roles[0]])
}
