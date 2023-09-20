package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Request(c *gin.Context) {
	c.Set("traceId", uuid.New().String())
}
