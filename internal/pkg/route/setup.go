package route

import (
	"github.com/Template7/backend/internal/pkg/route/handler"
	"github.com/Template7/backend/internal/pkg/route/middle_ware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.LoadHTMLFiles("resource/template/facebook_login.html")

	r.GET("", handler.HelloPage)

	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/test/graceful-shutdown", handler.TestGracefulShutdown)
	}

	apiV1 := r.Group("/api/v1")

	// user
	user := apiV1.Group("/users", middle_ware.AuthUserToken)
	user.GET("/:userId", handler.GetInfo)
	user.PUT("/:user-id", handler.UpdateUser)

	// sign up
	signUp := apiV1.Group("/sign-up")
	signUp.POST("/verification", handler.SendVerifyCode)
	signUp.POST("/confirmation", handler.ConfirmVerifyCode)

	// sign in
	signIn := apiV1.Group("/sign-in")
	signIn.POST("/mobile", handler.MobileSignIn)
	signIn.GET("/facebook/home", handler.FacebookSignInHome)
	signIn.GET("/facebook", handler.FacebookSignIn)
	signIn.GET("/facebook/callback", handler.FacebookSignInCallback)

	// admin
	adminV1 := r.Group("/admin/v1", middle_ware.AuthAdmin)

	adminV1.POST("/user", handler.CreateUser)
	//adminV1.DELETE("/user", handler.DeleteUser)

	// tokenless api

	// special case, skip auth token due to expired
	apiV1.PUT("/users/:user-id/token", handler.RefreshToken)
}
