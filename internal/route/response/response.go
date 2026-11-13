// Package response writes the JSON envelope shared by every endpoint and is
// the one place request errors are logged.
package response

import (
	"net/http"

	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

// HeaderRequestId carries the server-generated request id.
const HeaderRequestId = "request_id"

// RequestId returns the id assigned to the current request.
func RequestId(c *gin.Context) string {
	return c.GetHeader(HeaderRequestId)
}

// Ok returns the success envelope; embed it in endpoint-specific response types.
func Ok(c *gin.Context) types.HttpRespBase {
	return types.HttpRespBase{
		RequestId: RequestId(c),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	}
}

// JSON writes a 200 response with the given body.
func JSON(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, body)
}

// Error aborts the request with the status and code of err. Errors that are not
// *t7Error.Error become 500 Unknown. The error is logged here, once: 5xx as
// Error, 4xx as Info.
func Error(c *gin.Context, err error) {
	t7Err, ok := t7Error.ToT7Error(err)
	if !ok {
		t7Err = t7Error.Unknown
	}
	status := t7Err.GetStatus()

	log := logger.GetLogger().WithContext(c).
		WithService("http").
		With("method", c.Request.Method).
		With("path", c.FullPath()).
		With("status", status).
		With("code", t7Err.Code).
		WithError(err)
	if status >= http.StatusInternalServerError {
		log.Error("request failed")
	} else {
		log.Info("request rejected")
	}

	c.AbortWithStatusJSON(status, types.HttpRespBase{
		RequestId: RequestId(c),
		Code:      int(t7Err.Code),
		Message:   t7Err.Message,
	})
}

// BadRequest aborts with InvalidBody; detail is logged but not returned.
func BadRequest(c *gin.Context, detail string) {
	Error(c, t7Error.InvalidBody.WithDetail(detail))
}
