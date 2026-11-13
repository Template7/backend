package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service wallet.Service
	log     *logger.Logger
}

func NewWalletController(service wallet.Service, log *logger.Logger) *WalletController {
	return &WalletController{
		service: service,
		log:     log.WithService("walletController"),
	}
}

func toCurrency(name string) walletV1.Currency {
	return walletV1.Currency(walletV1.Currency_value[name])
}

// GetWallet
// @Summary Get wallet
// @Tags V1,Wallet
// @version 1.0
// @Success 200 {object} types.HttpGetWalletResp "Response"
// @failure 404 {object} types.HttpRespError "wallet not found"
// @Security BearerAuth
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId} [get]
func (w *WalletController) GetWallet(c *gin.Context) {
	userId, ok := currentUserId(c)
	if !ok {
		return
	}
	data, err := w.service.GetWallet(c, userId, c.Param("walletId"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.JSON(c, types.HttpGetWalletResp{
		HttpRespBase: response.Ok(c),
		Data: types.HttpGetUserWalletsRespData{
			Id:       data.GetId(),
			Balances: toBalances(data.GetBalances()),
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
	var req types.HttpWalletDepositReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	if err := w.service.Deposit(c, userId, c.Param("walletId"), toCurrency(req.Currency), req.Amount, req.Note); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, response.Ok(c))
}

// Withdraw
// @Summary Wallet withdraw
// @Tags V1,Wallet
// @version 1.0
// @Param request body types.HttpWalletWithdrawReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError "invalid body or insufficient balance"
// @Security BearerAuth
// @Param walletId path string true "Wallet ID"
// @Router /api/v1/wallets/{walletId}/withdraw [post]
func (w *WalletController) Withdraw(c *gin.Context) {
	var req types.HttpWalletWithdrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	if err := w.service.Withdraw(c, userId, c.Param("walletId"), toCurrency(req.Currency), req.Amount, req.Note); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, response.Ok(c))
}

// Transfer
// @Summary Wallet transfer
// @Tags V1,Wallet
// @version 1.0
// @Param request body types.HttpTransferMoneyReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError "invalid body or insufficient balance"
// @failure 404 {object} types.HttpRespError "source wallet not found or not owned by the caller"
// @Security BearerAuth
// @Router /api/v1/transfer [post]
func (w *WalletController) Transfer(c *gin.Context) {
	var req types.HttpTransferMoneyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	if err := w.service.Transfer(c, userId, req.FromWalletId, req.ToWalletId, toCurrency(req.Currency), req.Amount, req.Note); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, response.Ok(c))
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
	userId, ok := currentUserId(c)
	if !ok {
		return
	}
	history, err := w.service.GetBalanceHistory(c, userId, c.Param("walletId"))
	if err != nil {
		response.Error(c, err)
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

	response.JSON(c, types.HttpGetWalletBalanceHistoryResp{
		HttpRespBase: response.Ok(c),
		Data:         data,
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
	currency, ok := walletV1.Currency_value[c.Param("currency")]
	if !ok {
		response.BadRequest(c, "unsupported currency: "+c.Param("currency"))
		return
	}
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	history, err := w.service.GetBalanceHistoryByCurrency(c, userId, c.Param("walletId"), walletV1.Currency(currency))
	if err != nil {
		response.Error(c, err)
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

	response.JSON(c, types.HttpGetWalletBalanceHistoryByCurrencyResp{
		HttpRespBase: response.Ok(c),
		Data:         data,
	})
}

// GetUserWallets
// @Summary Get user wallets
// @Tags V1,User,Wallet
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpGetUserWalletsResp "Response"
// @failure 400 {object} types.HttpRespError
// @Security BearerAuth
// @Router /api/v1/user/wallets [get]
func (w *WalletController) GetUserWallets(c *gin.Context) {
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	uws, err := w.service.GetUserWallets(c, userId)
	if err != nil {
		response.Error(c, err)
		return
	}
	rd := make([]types.HttpGetUserWalletsRespData, len(uws))
	for i, uw := range uws {
		rd[i] = types.HttpGetUserWalletsRespData{
			Id:       uw.Id,
			Balances: toBalances(uw.Balances),
		}
	}

	response.JSON(c, types.HttpGetUserWalletsResp{
		HttpRespBase: response.Ok(c),
		Data:         rd,
	})
}
