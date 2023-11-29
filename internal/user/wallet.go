package user

import (
	"context"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *Service) GetUserWallets(ctx context.Context, userId string) (data []*v1.Wallet) {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user wallets")

	for _, uw := range s.db.GetUserWallets(ctx, userId) {
		data = append(data, uw.ToProto(ctx))
	}
	return data
}
