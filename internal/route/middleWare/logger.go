package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Request(c *gin.Context) {
	uId := uuid.New().String()
	c.Request.Header.Add(HeaderRequestId, uId)
	c.Set("traceId", uId)
}
