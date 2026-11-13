package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request assigns a server-generated request id, overriding any value sent by the client.
func (m *Controller) Request(c *gin.Context) {
	uId := uuid.New().String()
	c.Request.Header.Set(HeaderRequestId, uId)
	c.Set("traceId", uId)
}
