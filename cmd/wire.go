//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Template7/backend/internal"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/route/handler"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/backend/internal/wallet"
	"github.com/google/wire"
)

func InitializeApp() *App {
	wire.Build(
		NewApp,
		config.New,
		internal.ProvideLogger,
		internal.ProvideCacheCore,
		internal.ProvideSqlCore,
		internal.ProvideNoSqlCore,
		db.New,
		cache.New,
		handler.NewAuthController,
		handler.NewUserController,
		handler.NewWalletController,
		auth.New,
		user.New,
		wallet.New,
		middleware.New,
	)
	return &App{}
}
