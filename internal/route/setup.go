package route

import (
	"github.com/Template7/backend/internal/route/handler"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.Request, middleware.RecoverMiddleware)
	r.GET("", handler.HelloPage)

	apiV1 := r.Group("/api/v1")

	// login
	login := apiV1.Group("/login")
	login.POST("/native", handler.NativeLogin)

	apiV1.Use(middleware.AuthToken, middleware.Permission)

	// user
	user := apiV1.Group("/user")
	user.GET("/info", middleware.CheckAccountStatusInitialized, handler.GetUserInfo)
	user.GET("/wallets", middleware.CheckAccountStatusInitialized, handler.GetUserWallets)
	user.PUT("/info", middleware.CheckAccountStatusActivated, handler.UpdateUser)

	// wallet
	wallet := apiV1.Group("/wallets/:walletId", middleware.AuthUserWallet)
	wallet.GET("", middleware.CheckAccountStatusInitialized, handler.GetWallet)
	wallet.POST("/deposit", middleware.CheckAccountStatusActivated, handler.Deposit)
	wallet.POST("/withdraw", middleware.CheckAccountStatusActivated, handler.Withdraw)
	wallet.GET("/currencies/:currency/record", middleware.CheckAccountStatusActivated, handler.GetWalletBalanceRecord)

	// transfer
	transfer := apiV1.Group("/transfer", middleware.CheckAccountStatusActivated)
	transfer.POST("", handler.Transfer)

	// admin
	adminV1 := r.Group("/admin/v1", middleware.AuthToken, middleware.Permission)
	adminV1.POST("/user", handler.CreateUser, handler.HandleActivationCode)
	adminV1.DELETE("/users/:userId", handler.DeleteUser)
	adminV1.POST("/users/:userId/activate", handler.ActivateUser)

	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/test/graceful-shutdown", handler.TestGracefulShutdown)
	}
}
