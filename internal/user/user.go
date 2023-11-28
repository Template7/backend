package user

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/backend/internal/t7Error"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
)

func (s *Service) GetInfo(ctx context.Context, userId string) (*userV1.UserInfoResponse, error) {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user info")

	data, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		log.WithError(err).Error("fail to get user by id")
		return nil, t7Error.DbOperationFail.WithDetail(err.Error())
	}

	resp := userV1.UserInfoResponse{
		UserId:   data.Id,
		Role:     s.authSvc.GetUserRole(ctx, data.Id),
		Status:   data.Status,
		Nickname: data.Info.Nickname,
		Email:    data.Email,
	}
	log.With("userInfo", &resp).Debug("got user info")
	return &resp, nil
}

func (s *Service) UpdateInfo(ctx context.Context, userId string, info entity.UserInfo) error {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("update user info")

	if err := s.db.UpdateUserInfo(ctx, userId, info); err != nil {
		log.WithError(err).Error("fail to update user info")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}
	return nil
}
