package route

import (
	"github.com/Template7/backend/internal/pkg/route/handler"
	middleware "github.com/Template7/backend/internal/pkg/route/middleWare"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.Request, middleware.RecoverMiddleware)
	r.GET("", handler.HelloPage)

	r.Use(middleware.AuthToken, middleware.Permission)

	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/test/graceful-shutdown", handler.TestGracefulShutdown)
	}

	apiV1 := r.Group("/api/v1")

	// user
	user := apiV1.Group("/users", middleWare.AuthUserToken)
	user.GET("/:userId", handler.GetInfo)
	user.PUT("/:userId", handler.UpdateUser)

	// login
	login := apiV1.Group("/login")
	login.POST("/native", handler.NativeLogin)

	// wallet
	wallet := apiV1.Group("/wallet", middleWare.AuthUserToken, middleWare.AuthActiveUser)
	wallet.GET("/:walletId", handler.GetWallet)
	wallet.POST("/deposit", handler.Deposit)
	wallet.POST("/withdraw", handler.Withdraw)

	// transaction
	transaction := apiV1.Group("/transaction", middleWare.AuthUserToken, middleWare.AuthActiveUser)
	transaction.POST("", handler.MakeTransfer)

	// admin
	adminV1 := r.Group("/admin/v1", middleWare.AuthAdmin)
	adminV1.POST("/user", handler.CreateUser)
	//adminV1.DELETE("/user", handler.DeleteUser)

	// tokenless api

	// special case, skip auth token due to expired
	//apiV1.PUT("/users/:user-id/token", handler.RefreshToken)
}
