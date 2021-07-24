package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type hello struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// HelloPage
// @Summary Hello Page
// @Tags Hello
// @version 1.0
// @produce json
// @Success 200 {object} hello
// @Router / [get]
func HelloPage(c *gin.Context) {
	log.Debug("handle hello page")

	c.JSON(http.StatusOK, hello{
		Message:   "hello",
		Timestamp: time.Now(),
	})
	return
}

func TestGracefulShutdown(c *gin.Context) {
	log.Debug("handle test graceful shutdown")

	time.Sleep(10 * time.Second)

	c.String(http.StatusOK, "response after 10 second")
	return
}
