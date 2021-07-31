package handler

import (
	"backend/internal/pkg/admin"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// AdminSignIn
// @Summary Admin sign in
// @Tags v1,SignIn,Admin
// @version 1.0
// @Param smsRequest body collection.Admin true "Admin object"
// @produce json
// @Success 200 {object} collection.Token "Token object"
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /admin/v1/sign-in [post]
func AdminSignIn(c *gin.Context) {
	log.Debug("handle admin sign in")

	var data collection.Admin
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}
	token, err := admin.SignIn(data)
	if err != nil {
		log.Error("fail to sign in admin: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, token)
	return
}
