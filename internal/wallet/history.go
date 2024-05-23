package wallet

import (
	"context"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

func (s *Service) GetBalanceHistoryByCurrency(ctx context.Context, walletId string, currency walletV1.Currency) []walletV1.CurrencyBalanceRecord {
	// TODO
	return nil
}
