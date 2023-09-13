package route

import (
	"github.com/Template7/backend/internal/route/handler"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.Request, middleware.RecoverMiddleware)
	r.GET("", handler.HelloPage)

	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/test/graceful-shutdown", handler.TestGracefulShutdown)
	}

	apiV1 := r.Group("/api/v1")

	// login
	login := apiV1.Group("/login")
	login.POST("/native", handler.NativeLogin)

	apiV1.Use(middleware.AuthToken, middleware.Permission)

	// user
	user := apiV1.Group("/users")
	user.GET("/:userId/info", handler.GetUserInfo)
	user.PUT("/:userId", handler.UpdateUser)

	// wallet
	wallet := apiV1.Group("/wallet")
	wallet.GET("/:walletId", handler.GetWallet)
	wallet.POST("/deposit", handler.Deposit)
	wallet.POST("/withdraw", handler.Withdraw)

	// transaction
	transaction := apiV1.Group("/transaction")
	transaction.POST("", handler.MakeTransfer)

	// admin
	adminV1 := r.Group("/admin/v1")
	adminV1.POST("/user", handler.CreateUser)
	//adminV1.DELETE("/user", handler.DeleteUser)
}
