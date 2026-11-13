package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *service) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (string, error) {
	log := s.log.WithContext(ctx).With("data", req.Username)
	log.Debug("create user")

	hp, err := hashedPassword(req.Password)
	if err != nil {
		log.WithError(err).Error("fail to hash password")
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

	data := models.User{
		Id:       userId,
		Username: req.Username,
		Password: hp,
		Info: models.UserInfo{
			Nickname: req.Nickname,
		},
		Email:  req.Email,
		Status: authV1.AccountStatus_initialized,
	}
	if err := s.db.CreateUser(ctx, data); err != nil {
		log.WithError(err).Error("fail to create user")
		// undo the role so no orphan casbin entry is left behind
		if _, rErr := s.core.DeleteRolesForUser(userId); rErr != nil {
			log.WithError(rErr).Error("fail to roll back user role")
		}
		return "", t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("user created")
	return userId, nil
}

func (s *service) GenActivationCode(ctx context.Context, userId string) (string, error) {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("gen user activation code")

	actCode := uuid.NewString()
	if err := s.cache.SetUserActivationCode(ctx, userId, actCode); err != nil {
		log.WithError(err).Error("fail to set user activation code")
		return "", t7Error.RedisOperationFail.WithDetail(err.Error())
	}

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

	if actCode == "" || subtle.ConstantTimeCompare([]byte(code), []byte(actCode)) != 1 {
		log.Warn("invalid activation code")
		return false
	}

	if err := s.db.SetUserStatus(ctx, userId, authV1.AccountStatus_activated); err != nil {
		log.WithError(err).Error("fail to update user status")
		return false
	}

	// one-time code
	if err := s.cache.DeleteUserActivationCode(ctx, userId); err != nil {
		log.WithError(err).Warn("fail to delete used activation code")
	}

	log.Debug("user activated")
	return true
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("delete user")

	// delete the data first: if it fails, the user keeps working as before
	if err := s.db.DeleteUser(ctx, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("user not found")
			return t7Error.UserNotfound
		}
		log.WithError(err).Error("fail to delete user")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	// the user row is gone, so a leftover role can't be used to log in; just log it
	if _, err := s.core.DeleteUser(userId); err != nil {
		log.WithError(err).Error("fail to delete user from casbin")
	}
	if err := s.cache.DeleteUserActivationCode(ctx, userId); err != nil {
		log.WithError(err).Warn("fail to delete activation code")
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
