package user

import (
	"context"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *Service) GetUserWallets(ctx context.Context, userId string) (data []*v1.Wallet) {
	log := s.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user wallets")

	uws := s.db.GetUserWallets(ctx, userId)
	if len(uws) == 0 {
		log.Warn("user has no wallets")
		return
	}

	// wallet entity to proto wallet
	for _, uw := range s.db.GetUserWallets(ctx, userId) {
		var blc []*v1.Balance
		for _, b := range uw.Balance {
			if _, ok := v1.Currency_value[b.Currency]; !ok {
				log.With("currency", b.Currency).Warn("unsupported currency")
				continue
			}
			blc = append(blc, &v1.Balance{
				Currency: v1.Currency(v1.Currency_value[b.Currency]),
				Amount:   b.Amount.String(),
			})
		}
		data = append(data, &v1.Wallet{
			Id:       uw.Id,
			Balances: blc,
		})
	}
	return data
}
