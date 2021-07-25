package handler

import (
	"backend/internal/pkg/admin"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

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
