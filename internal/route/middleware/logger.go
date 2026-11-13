package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// traceIdKey is the context key the shared logger reads (it uses a plain string).
const traceIdKey = "traceId"

// Request assigns a server-generated request id, overriding any value sent by the
// client, and stores it on both the gin context and the request context, so logs
// from services (which get c.Request.Context()) carry the same trace id.
func (m *Controller) Request(c *gin.Context) {
	uId := uuid.New().String()
	c.Request.Header.Set(HeaderRequestId, uId)
	c.Set(traceIdKey, uId)
	//nolint:staticcheck // SA1029: the shared logger looks the trace id up by this string key
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), traceIdKey, uId))
}
