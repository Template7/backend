package user

import (
	"context"
	"crypto/subtle"

	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/google/uuid"
)

func (s *service) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (string, string, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return "", "", t7Error.HashFail.WithDetail(err.Error())
	}

	userId := uuid.NewString()
	user := models.User{
		Id:       userId,
		Username: req.Username,
		Password: hash,
		Info:     models.UserInfo{Nickname: req.Nickname},
		Email:    req.Email,
		Status:   authV1.AccountStatus_initialized,
	}
	wallet := models.Wallet{Id: uuid.NewString(), UserId: userId}

	err = s.db.Transaction(ctx, func(tx db.Client) error {
		if err := tx.CreateUser(ctx, user); err != nil {
			return err
		}
		return tx.CreateWallet(ctx, wallet, defaultCurrencies)
	})
	if err != nil {
		if _, getErr := s.db.GetUser(ctx, req.Username); getErr == nil {
			return "", "", t7Error.UserAlreadyExist
		}
		return "", "", t7Error.DbOperationFail.WithDetail(err.Error())
	}

	// the role lives in casbin, outside the transaction: undo the user if it fails
	if err := s.auth.AddUserRole(ctx, userId, req.Role); err != nil {
		if delErr := s.db.DeleteUser(ctx, userId); delErr != nil {
			s.log.WithContext(ctx).WithError(delErr).With("userId", userId).Error("fail to roll back user after role error")
		}
		return "", "", err
	}

	code := uuid.NewString()
	if err := s.cache.SetUserActivationCode(ctx, userId, code); err != nil {
		// the account exists; an admin can't activate it until a new code is issued
		return userId, "", t7Error.RedisOperationFail.WithDetail(err.Error())
	}
	return userId, code, nil
}

func (s *service) ActivateUser(ctx context.Context, userId string, code string) error {
	stored, err := s.cache.GetUserActivationCode(ctx, userId)
	if err != nil {
		// missing or expired code
		return t7Error.IncorrectVerifyCode.WithDetail(err.Error())
	}
	if code == "" || subtle.ConstantTimeCompare([]byte(stored), []byte(code)) != 1 {
		return t7Error.IncorrectVerifyCode
	}

	if err := s.db.SetUserStatus(ctx, userId, authV1.AccountStatus_activated); err != nil {
		if db.IsNotFound(err) {
			return t7Error.UserNotfound
		}
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	// one-time code; the account is already active, so only log a failure
	if err := s.cache.DeleteUserActivationCode(ctx, userId); err != nil {
		s.log.WithContext(ctx).WithError(err).With("userId", userId).Warn("fail to delete used activation code")
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	// data first: if it fails, the account keeps working as before
	if err := s.db.DeleteUser(ctx, userId); err != nil {
		if db.IsNotFound(err) {
			return t7Error.UserNotfound
		}
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	// the user row is gone, so leftovers can't be used; only log failures
	log := s.log.WithContext(ctx).With("userId", userId)
	if err := s.auth.DeleteUserRoles(ctx, userId); err != nil {
		log.WithError(err).Error("fail to delete user roles")
	}
	if err := s.cache.DeleteUserActivationCode(ctx, userId); err != nil {
		log.WithError(err).Warn("fail to delete activation code")
	}
	return nil
}

func (s *service) GetInfo(ctx context.Context, userId string) (*userV1.UserInfoResponse, error) {
	data, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	role, err := s.auth.GetUserRole(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &userV1.UserInfoResponse{
		UserId:   data.Id,
		Role:     role,
		Status:   data.Status,
		Nickname: data.Info.Nickname,
		Email:    data.Email,
	}, nil
}

func (s *service) UpdateInfo(ctx context.Context, userId string, info models.UserInfo) error {
	if err := s.db.UpdateUserInfo(ctx, userId, info); err != nil {
		if db.IsNotFound(err) {
			return t7Error.UserNotfound
		}
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}
	return nil
}

func (s *service) GetStatus(ctx context.Context, userId string) (authV1.AccountStatus, error) {
	data, err := s.getUser(ctx, userId)
	if err != nil {
		return 0, err
	}
	return data.Status, nil
}

func (s *service) getUser(ctx context.Context, userId string) (models.User, error) {
	data, err := s.db.GetUserById(ctx, userId)
	if db.IsNotFound(err) {
		return models.User{}, t7Error.UserNotfound
	}
	if err != nil {
		return models.User{}, t7Error.DbOperationFail.WithDetail(err.Error())
	}
	return data, nil
}
