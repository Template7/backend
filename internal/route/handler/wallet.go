package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/db"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

type WalletController struct {
	service *wallet.Service
	log     *logger.Logger
}

func NewWalletController(service *wallet.Service, log *logger.Logger) *WalletController {
	return &WalletController{
		service: service,
		log:     log.With("userService", "walletController"),
	}
}

// GetWallet
// @Summary Get wallet
// @Tags V1,Wallet
// @version 1.0
// @Success 200 {object} types.HttpGetWalletResp "Response"
// @failure 400 {object} types.HttpRespError
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId} [get]
func (w *WalletController) GetWallet(c *gin.Context) {
	log := w.log.WithContext(c)
	log.Debug("handle get wallet")

	data, err := w.service.GetWallet(c, c.Param("walletId"))
	if err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to get wallet")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	bls := make([]types.HttpGetUserWalletsRespDataBalance, len(data.Balances))
	for i, bl := range data.GetBalances() {
		bls[i] = types.HttpGetUserWalletsRespDataBalance{
			Currency: bl.GetCurrency().String(),
			Amount:   bl.GetAmount(),
		}
	}

	c.JSON(http.StatusOK, types.HttpGetWalletResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: types.HttpGetUserWalletsRespData{
			Id:       data.Id,
			Balances: bls,
		},
	})
}

// Deposit
// @Summary Wallet deposit
// @Tags V1,Wallet
// @version 1.0
// @Param request body types.HttpWalletDepositReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId}/deposit [post]
func (w *WalletController) Deposit(c *gin.Context) {
	log := w.log.WithContext(c)
	log.Debug("handle deposit")

	var req types.HttpWalletDepositReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	if err := w.service.Deposit(c, c.Param("walletId"), v1.Currency(v1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to deposit")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpRespBase{
		RequestId: c.GetHeader(middleware.HeaderRequestId),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	})
}

// Withdraw
// @Summary Wallet withdraw
// @Tags V1,Wallet
// @version 1.0
// @Param request body types.HttpWalletWithdrawReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId}/withdraw [post]
func (w *WalletController) Withdraw(c *gin.Context) {
	log := w.log.WithContext(c)
	log.Debug("handle withdraw")

	var req types.HttpWalletWithdrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	if err := w.service.Withdraw(c, c.Param("walletId"), v1.Currency(v1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to withdraw")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpRespBase{
		RequestId: c.GetHeader(middleware.HeaderRequestId),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	})
}

// Transfer
// @Summary Wallet transfer
// @Tags V1,Wallet
// @version 1.0
// @Param request body types.HttpTransferMoneyReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Router /api/v1/transfer [post]
func (w *WalletController) Transfer(c *gin.Context) {
	log := w.log.WithContext(c)
	log.WithContext(c).Debug("handle make transfer")

	var req types.HttpTransferMoneyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	if err := w.service.Transfer(c, req.FromWalletId, req.ToWalletId, v1.Currency(v1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to transfer")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpRespBase{
		RequestId: c.GetHeader(middleware.HeaderRequestId),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	})
}

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
func (w *WalletController) GetWalletBalanceRecord(c *gin.Context) {
	log := w.log.WithContext(c)
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

	records, err := db.New().GetWalletBalanceHistory(c, wId, cur)
	if err != nil {
		log.WithError(err).Error("fail to get wallet balance history")
		c.JSON(http.StatusInternalServerError, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.DbOperationFail.Code),
			Message:   t7Error.DbOperationFail.Message,
		})
		return
	}

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
			Amount:       r.Amount,
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
