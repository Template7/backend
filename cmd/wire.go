//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/route/handler"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	"github.com/google/wire"
)

func InitializeApp() *App {
	wire.Build(
		NewApp,
		handler.NewAuthController,
		handler.NewUserController,
		handler.NewWalletController,
		auth.New,
		user.New,
		wallet.New,
		config.New,
		//db.New,
		logger.New,
	)
	return &App{}
}
