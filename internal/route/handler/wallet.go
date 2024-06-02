package handler

import (
	"github.com/Template7/backend/api/types"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @Security BearerAuth
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
// @Security BearerAuth
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

	if err := w.service.Deposit(c, c.Param("walletId"), walletV1.Currency(walletV1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
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
// @Security BearerAuth
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

	if err := w.service.Withdraw(c, c.Param("walletId"), walletV1.Currency(walletV1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
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
// @Security BearerAuth
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

	if err := w.service.Transfer(c, req.FromWalletId, req.ToWalletId, walletV1.Currency(walletV1.Currency_value[req.Currency]), req.Amount, req.Note); err != nil {
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

// GetWalletBalanceHistory
// @Summary Get wallet balance record
// @Tags V1,Wallet
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpGetWalletBalanceHistoryResp "Response"
// @failure 400 {object} types.HttpRespError
// @Security BearerAuth
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId}/history [get]
func (w *WalletController) GetWalletBalanceHistory(c *gin.Context) {
	log := w.log.WithContext(c)
	log.Debug("handle get wallet balance history")

	wId := c.Param("walletId")
	if wId == "" {
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	history := w.service.GetBalanceHistory(c, wId)

	if history == nil {
		c.JSON(http.StatusInternalServerError, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.DbOperationFail.Code),
			Message:   t7Error.DbOperationFail.Message,
		})
		return
	}

	data := make([]types.HttpGetWalletBalanceHistoryData, len(history))
	for i, h := range history {
		data[i] = types.HttpGetWalletBalanceHistoryData{
			Direction:     h.Direction.String(),
			Currency:      h.Currency.String(),
			Amount:        h.Amount,
			BalanceBefore: h.BalanceBefore,
			BalanceAfter:  h.BalanceAfter,
			Timestamp:     h.Timestamp,
			Note:          h.Note,
		}
	}

	c.JSON(http.StatusOK, types.HttpGetWalletBalanceHistoryResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: data,
	})
}

// GetWalletBalanceHistoryByCurrency
// @Summary Get wallet balance record
// @Tags V1,Wallet
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpGetWalletBalanceHistoryByCurrencyResp "Response"
// @failure 400 {object} types.HttpRespError
// @Security BearerAuth
// @Param walletId path string true "Wallet ID"
// @Param currency path string true "Currency" Enums(ntd, cny, usd, jpy)
// @Router /api/v1/wallets/{walletId}/currencies/{currency}/history [get]
func (w *WalletController) GetWalletBalanceHistoryByCurrency(c *gin.Context) {
	log := w.log.WithContext(c)
	log.Debug("handle get wallet balance history by currency")

	wId := c.Param("walletId")
	if wId == "" {
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	currency, ok := walletV1.Currency_value[c.Param("currency")]
	if !ok {
		log.With("currency", c.Param("currency")).Info("unsupported currency")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	history := w.service.GetBalanceHistoryByCurrency(c, wId, walletV1.Currency(currency))

	if history == nil {
		c.JSON(http.StatusInternalServerError, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.DbOperationFail.Code),
			Message:   t7Error.DbOperationFail.Message,
		})
		return
	}

	data := make([]types.HttpGetWalletBalanceHistoryByCurrencyData, len(history))
	for i, h := range history {
		data[i] = types.HttpGetWalletBalanceHistoryByCurrencyData{
			Direction:     h.Direction.String(),
			Amount:        h.Amount,
			BalanceBefore: h.BalanceBefore,
			BalanceAfter:  h.BalanceAfter,
			Timestamp:     h.Timestamp,
			Note:          h.Note,
		}
	}

	c.JSON(http.StatusOK, types.HttpGetWalletBalanceHistoryByCurrencyResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: data,
	})
}
