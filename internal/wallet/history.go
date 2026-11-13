package wallet

import (
	"context"

	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *service) GetBalanceHistory(ctx context.Context, userId string, walletId string) ([]*walletV1.BalanceHistory, error) {
	if err := s.checkOwner(ctx, userId, walletId); err != nil {
		return nil, err
	}
	history, err := s.db.GetWalletBalanceHistory(ctx, walletId)
	if err != nil {
		return nil, dbError(err)
	}

	data := make([]*walletV1.BalanceHistory, len(history))
	for i, h := range history {
		data[i] = &walletV1.BalanceHistory{
			Id:            h.Id,
			Direction:     walletV1.Direction(walletV1.Direction_value[h.Direction]),
			Currency:      walletV1.Currency(walletV1.Currency_value[h.Currency]),
			Amount:        h.Amount.String(),
			BalanceBefore: h.BalanceBefore.String(),
			BalanceAfter:  h.BalanceAfter.String(),
			Timestamp:     h.CreatedAt.UnixMilli(),
			Note:          h.Note,
		}
	}
	return data, nil
}

func (s *service) GetBalanceHistoryByCurrency(ctx context.Context, userId string, walletId string, currency walletV1.Currency) ([]*walletV1.CurrencyBalanceHistory, error) {
	if err := s.checkOwner(ctx, userId, walletId); err != nil {
		return nil, err
	}
	history, err := s.db.GetWalletBalanceHistoryByCurrency(ctx, walletId, currency.String())
	if err != nil {
		return nil, dbError(err)
	}

	data := make([]*walletV1.CurrencyBalanceHistory, len(history))
	for i, h := range history {
		data[i] = &walletV1.CurrencyBalanceHistory{
			Id:            h.Id,
			Direction:     walletV1.Direction(walletV1.Direction_value[h.Direction]),
			Amount:        h.Amount.String(),
			BalanceBefore: h.BalanceBefore.String(),
			BalanceAfter:  h.BalanceAfter.String(),
			Timestamp:     h.CreatedAt.UnixMilli(),
			Note:          h.Note,
		}
	}
	return data, nil
}
