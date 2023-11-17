package handler

import (
	"github.com/Template7/backend/api/types"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

// GetWallet
// @Summary Get wallet
// @Tags v1,wallet
// @version 1.0
// @Success 200 {object} types.HttpGetWalletResp "Response"
// @failure 400 {object} types.HttpRespError
// @Param walletId path string true "wallet id"
// @Router /api/v1/wallets/{walletId} [get]
func GetWallet(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get wallet")

	data, err := wallet.New().GetWallet(c, c.Param("walletId"))
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

func Deposit(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle deposit")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req v1.DepositRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	if err := wallet.New().Deposit(c, c.Param("walletId"), req.GetCurrency(), req.GetAmount()); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to deposit")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func Withdraw(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle withdraw")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req v1.WithdrawRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	if err := wallet.New().Withdraw(c, c.Param("walletId"), req.GetCurrency(), req.GetAmount()); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to withdraw")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func Transfer(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.WithContext(c).Debug("handle make transfer")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req v1.TransferRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	if err := wallet.New().Transfer(c, req.GetFromWalletId(), req.GetToWalletId(), req.GetCurrency(), req.GetAmount()); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to transfer")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
