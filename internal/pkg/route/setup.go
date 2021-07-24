package route

import (
	"backend/internal/pkg/route/handler"
	"backend/internal/pkg/route/middle_ware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.GET("", handler.HelloPage)

	if gin.Mode() == gin.DebugMode {
		//url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/index.html", config.New().Port))
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiV1 := r.Group("/api/v1")

	// user
	user := apiV1.Group("/users", middle_ware.AuthUser)

	user.GET("/:user-id", handler.GetInfo)
	user.PUT("/:user-id", handler.UpdateUser)
	//user.DELETE("/:user-id", handler.DeleteUser)
	user.PUT("/:user-id/token", handler.RefreshToken)
	//user.PUT("/:user-id/login-client", handler.UpdateLoginClient)
	//user.PUT("/:user-id/logout", handler.SignOut)

	// sign up
	signUp := apiV1.Group("/sign-up")
	signUp.POST("/verification", handler.SendVerifyCode)
	signUp.POST("/confirmation", handler.ConfirmVerifyCode)

	signIn := apiV1.Group("/sign-in")
	signIn.POST("/mobile/verification", handler.MobileSignIn)
	signIn.POST("/mobile/confirmation", handler.MobileSignInConfirm)
	//signIn.POST("/facebook", handler.FacebookSignIn)


	// admin
	adminV1 := r.Group("/admin/v1", middle_ware.AuthAdmin)
	adminV1.POST("/user", handler.CreateUser)
	adminV1.DELETE("/user", handler.DeleteUser)


}
