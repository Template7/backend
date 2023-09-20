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
	user.GET("/info", handler.GetUserInfo)
	user.GET("/wallets", handler.GetUserWallets)
	user.PUT("/info", handler.UpdateUser)
	user.POST("/new", handler.CreateUser)

	// wallet
	wallet := apiV1.Group("/wallets/:walletId", middleware.AuthUserWallet)
	wallet.GET("", handler.GetWallet)
	wallet.POST("/deposit", handler.Deposit)
	wallet.POST("/withdraw", handler.Withdraw)

	// transaction
	transaction := apiV1.Group("/transaction")
	transaction.POST("", handler.MakeTransfer)

	// admin
	adminV1 := r.Group("/admin/v1")
	adminV1.POST("/user", handler.CreateUser)
	//adminV1.DELETE("/user", handler.DeleteUser)

	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/test/graceful-shutdown", handler.TestGracefulShutdown)
	}
}
