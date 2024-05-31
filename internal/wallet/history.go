package wallet

import (
	"context"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *Service) GetBalanceHistoryByCurrency(ctx context.Context, walletId string, currency walletV1.Currency) []walletV1.CurrencyBalanceHistory {
	log := s.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get wallet balance history by currency")

	history, err := s.db.GetWalletBalanceHistoryByCurrency(ctx, walletId, currency.String())
	if err != nil {
		log.WithError(err).Error("fail to get balance history by currency")
		return nil
	}

	data := make([]walletV1.CurrencyBalanceHistory, len(history))
	for i, h := range history {
		data[i] = walletV1.CurrencyBalanceHistory{
			Id:            h.Id,
			Direction:     walletV1.Direction(walletV1.Direction_value[h.Direction]),
			Amount:        h.Amount.String(),
			BalanceBefore: h.BalanceBefore.String(),
			BalanceAfter:  h.BalanceAfter.String(),
			Timestamp:     h.CreatedAt.UnixMilli(),
			Note:          h.Note,
		}
	}
	return data
}

func (s *Service) GetBalanceHistory(ctx context.Context, walletId string) []*walletV1.BalanceHistory {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet balance history")

	history, err := s.db.GetWalletBalanceHistory(ctx, walletId)
	if err != nil {
		log.WithError(err).Error("fail to get balance history")
		return nil
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
	return data
}
