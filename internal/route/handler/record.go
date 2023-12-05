package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/db"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

// GetWalletBalanceRecord
// @Summary Get wallet balance record
// @Tags V1,Wallet
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpGetWalletBalanceRecordResp "Response"
// @failure 400 {object} types.HttpRespError
// @Param walletId path string true "Wallet ID"
// @Param currency path string true "Currency"
// @Router /api/v1/wallets/{walletId}/currencies/{currency}/record [get]
func GetWalletBalanceRecord(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get wallet balance record")

	wId := c.Param("walletId")
	cur := c.Param("currency")
	if wId == "" || cur == "" {
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	records := db.New().GetWalletBalanceHistory(c, wId, cur)
	if records == nil {
		c.JSON(http.StatusInternalServerError, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.DbOperationFail.Code),
			Message:   t7Error.DbOperationFail.Message,
		})
		return
	}

	data := make([]types.HttpGetWalletBalanceRecordRespData, len(records))
	for i, r := range records {
		data[i] = types.HttpGetWalletBalanceRecordRespData{
			RecordId:     r.RecordId,
			Io:           r.Io,
			AmountBefore: r.AmountBefore,
			AmountAfter:  r.AmountAfter,
			Timestamp:    r.Timestamp,
			Note:         r.Note,
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Before(data[j].Timestamp)
	})

	c.JSON(http.StatusOK, types.HttpGetWalletBalanceRecordResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: data,
	})

}
