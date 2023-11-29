package user

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *Service) GetUserWallets(ctx context.Context, userId string) (data []*v1.Wallet) {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user wallets")

	uws := s.db.GetUserWallets(ctx, userId)
	gbw := func() map[string][]entity.UserWalletBalance {
		r := map[string][]entity.UserWalletBalance{}
		for _, uw := range uws {
			if _, ok := r[uw.WalletId]; !ok {
				r[uw.WalletId] = []entity.UserWalletBalance{uw}
			} else {
				r[uw.WalletId] = append(r[uw.WalletId], uw)
			}
		}
		return r
	}()

	for wId, uws := range gbw {
		var bls []*v1.Balance
		for _, bl := range uws {
			bls = append(bls, &v1.Balance{
				Currency: v1.Currency(v1.Currency_value[bl.Currency]),
				Amount:   bl.Amount.String(),
			})
		}
		data = append(data, &v1.Wallet{
			Id:       wId,
			Balances: bls,
		})
	}
	return data
}
