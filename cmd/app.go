package main

import (
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/route/handler"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type App struct {
	route            *gin.Engine
	authController   *handler.AuthController
	userController   *handler.UserController
	walletController *handler.WalletController
	config           *config.Config
	Log              *logger.Logger
}

func (a *App) SetupRoutes() *gin.Engine {
	a.route = gin.New()
	a.route.Use(middleware.Request, middleware.RecoverMiddleware)

	a.route.GET("", handler.HelloPage)

	apiV1 := a.route.Group("/api/v1")

	// login
	login := apiV1.Group("/login")
	login.POST("/native", a.authController.NativeLogin)

	apiV1.Use(middleware.AuthToken, middleware.Permission)

	// user
	user := apiV1.Group("/user")
	user.GET("/info", middleware.CheckAccountStatusInitialized, a.userController.GetUserInfo)
	user.GET("/wallets", middleware.CheckAccountStatusInitialized, a.userController.GetUserWallets)
	user.PUT("/info", middleware.CheckAccountStatusActivated, a.userController.UpdateUser)

	// wallet
	wallet := apiV1.Group("/wallets/:walletId", middleware.AuthUserWallet)
	wallet.GET("", middleware.CheckAccountStatusInitialized, a.walletController.GetWallet)
	wallet.POST("/deposit", middleware.CheckAccountStatusActivated, a.walletController.Deposit)
	wallet.POST("/withdraw", middleware.CheckAccountStatusActivated, a.walletController.Withdraw)
	wallet.GET("/currencies/:currency/record", middleware.CheckAccountStatusActivated, a.walletController.GetWalletBalanceRecord)

	// transfer
	transfer := apiV1.Group("/transfer", middleware.CheckAccountStatusActivated)
	transfer.POST("", a.walletController.Transfer)

	// admin
	adminV1 := a.route.Group("/admin/v1", middleware.AuthToken, middleware.Permission)
	adminV1.POST("/user", a.userController.CreateUser, a.userController.HandleActivationCode)
	adminV1.DELETE("/users/:userId", a.userController.DeleteUser)
	adminV1.POST("/users/:userId/activate", a.userController.ActivateUser)

	if a.config.Env == "dev" {
		a.route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return a.route
}

func NewApp(authController *handler.AuthController, userController *handler.UserController, walletController *handler.WalletController, config *config.Config, log *logger.Logger) *App {
	return &App{
		authController:   authController,
		userController:   userController,
		walletController: walletController,
		config:           config,
		Log:              log.With("service", "app"),
	}
}
