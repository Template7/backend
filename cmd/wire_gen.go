// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/Template7/backend/internal"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/route/handler"
	"github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/backend/internal/wallet"
)

import (
	_ "github.com/Template7/backend/docs"
)

// Injectors from wire.go:

func InitializeApp() *App {
	configConfig := config.New()
	gormDB := internal.ProvideSqlCore(configConfig)
	client := internal.ProvideNoSqlCore(configConfig)
	logger := internal.ProvideLogger(configConfig)
	dbClient := db.New(gormDB, client, logger)
	redisClient := internal.ProvideCacheCore(configConfig)
	cacheInterface := cache.New(redisClient, logger)
	authAuth := auth.New(dbClient, gormDB, cacheInterface, logger, configConfig)
	authController := handler.NewAuthController(authAuth, logger)
	service := user.New(authAuth, dbClient, logger)
	userController := handler.NewUserController(service, authAuth, logger)
	walletService := wallet.New(dbClient, logger)
	walletController := handler.NewWalletController(walletService, logger)
	controller := middleware.New(service, authAuth, logger)
	app := NewApp(authController, userController, walletController, controller, configConfig, logger)
	return app
}
