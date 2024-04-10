package auth

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/backend/internal/t7Error"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/google/uuid"
)

func (s *service) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (string, error) {
	log := s.log.WithContext(ctx).With("data", req.Username)
	log.Debug("create user")

	hp, err := hashedPassword(req.Password)
	if err != nil {
		log.WithError(err).With("password", req.Password).Error("fail to hash password")
		return "", t7Error.DecodeFail.WithDetail(err.Error())
	}

	userId := uuid.NewString()
	ok, err := s.core.AddRoleForUser(userId, req.Role.String())
	if err != nil {
		log.WithError(err).Error("fail to add role for user")
		return "", t7Error.Unknown.WithDetail(err.Error())
	}
	if !ok {
		log.Warn("user already exist")
		return "", t7Error.UserAlreadyExist
	}

	data := entity.User{
		Id:       userId,
		Username: req.Username,
		Password: hp,
		Info: entity.UserInfo{
			Nickname: req.Nickname,
		},
		Email:  req.Email,
		Status: authV1.AccountStatus_initialized,
	}
	if err := s.db.CreateUser(ctx, data); err != nil {
		log.WithError(err).Error("fail to create user")
		return "", t7Error.DbOperationFail.WithDetail(err.Error())
	}

	actCode := uuid.NewString()
	if err := s.cache.SetUserActivationCode(ctx, userId, actCode); err != nil {
		log.WithError(err).Error("fail to get user activation code")
		return "", t7Error.RedisOperationFail.WithDetail(err.Error())
	}

	log.Debug("user created")
	return actCode, nil
}

func (s *service) ActivateUser(ctx context.Context, userId string, actCode string) bool {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("activate user")

	code, err := s.cache.GetUserActivationCode(ctx, userId)
	if err != nil {
		log.WithError(err).Error("fail to get user activation code")
		return false
	}

	if code != actCode {
		log.Warn("invalid activation code")
		return false
	}

	log.Debug("user activated")
	return true
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("delete user")

	_, err := s.core.DeleteUser(userId)
	if err != nil {
		log.WithError(err).Error("fail to delete user from casbin api")
	}

	if err := s.db.DeleteUser(ctx, userId); err != nil {
		log.WithError(err).Error("fail to delete user")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	return nil
}

func (s *service) GetUserRole(ctx context.Context, userId string) authV1.Role {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user role")

	roles, err := s.core.GetRolesForUser(userId)
	if err != nil {
		log.WithError(err).Error("unable to get user role")
		return -1
	}
	if len(roles) == 0 {
		log.Warn("user has no roles")
		return -1
	}
	if len(roles) > 1 {
		log.With("roles", roles).Warn("user has multiple roles")
	}

	log.With("role", roles).Debug("got user role")
	return authV1.Role(authV1.Role_value[roles[0]])
}

func (s *service) GetUserStatus(ctx context.Context, userId string) authV1.AccountStatus {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user status")

	data, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		log.WithError(err).Error("fail to get user")
		return -1
	}

	return data.Status
}

func (s *service) ActiveUser(ctx context.Context, userId string) error {
	// TODO
	return nil
}
