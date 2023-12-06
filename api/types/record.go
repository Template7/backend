package types

import (
	"github.com/shopspring/decimal"
	"time"
)

type HttpGetWalletBalanceRecordResp struct {
	HttpRespBase
	Data []HttpGetWalletBalanceRecordRespData
}

type HttpGetWalletBalanceRecordRespData struct {
	RecordId     int64           `json:"recordId"`
	Io           string          `json:"io"` // in/out
	Amount       decimal.Decimal `json:"amount"`
	AmountBefore decimal.Decimal `json:"amountBefore"`
	AmountAfter  decimal.Decimal `json:"amountAfter"`
	Timestamp    time.Time       `json:"timestamp"` // record created_at
	Note         string          `json:"note"`
}
