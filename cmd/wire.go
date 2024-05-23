//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/route/handler"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/backend/internal/wallet"
	commonDb "github.com/Template7/common/db"
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
		middleware.New,
		db.New,
		commonDb.NewSql,
		cache.New,
		logger.New,
	)
	return &App{}
}
