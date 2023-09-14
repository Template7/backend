package auth

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/backend/internal/t7Error"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/google/uuid"
)

func (s *service) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) error {
	log := s.log.WithContext(ctx).With("data", req.Username)
	log.Debug("create user")

	hp, err := hashedPassword(req.Password)
	if err != nil {
		log.WithError(err).With("password", req.Password).Error("fail to hash password")
		return t7Error.DecodeFail.WithDetail(err.Error())
	}

	userId := uuid.New()
	ok, err := s.core.AddRoleForUser(userId.String(), req.Role.String())
	if err != nil {
		log.WithError(err).Error("fail to add role for user")
		return t7Error.Unknown.WithDetail(err.Error())
	}
	if !ok {
		log.Warn("user already exist")
		return t7Error.UserAlreadyExist
	}
	data := entity.User{
		Id:       userId,
		Username: req.Username,
		Password: hp,
		Info: entity.UserInfo{
			NickName: req.Nickname,
		},
		Email:  req.Email,
		Status: authV1.AccountStatus_Initialized,
	}
	if err := s.db.CreateUser(ctx, data); err != nil {
		log.WithError(err).Error("fail to create user")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("user created")
	return nil
}

func (s *service) GetUserRole(ctx context.Context, username string) authV1.Role {
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
	return authV1.Role(authV1.Role_value[roles[0]])
}
