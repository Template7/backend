package middleware

import (
	"errors"
	"fmt"
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// RecoverMiddleware turns a panic in any later handler into a 500 response
// instead of a dropped connection.
func (m *Controller) RecoverMiddleware(c *gin.Context) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		// panic values are not always errors, e.g. panic("msg")
		err, isErr := rec.(error)
		if isErr && isBrokenPipe(err) {
			// client went away; nothing to respond to
			_ = c.Error(err)
			c.Abort()
			return
		}

		m.log.WithContext(c).
			With("panic", fmt.Sprint(rec)).
			With("method", c.Request.Method).
			With("path", c.Request.URL.Path).
			With("stack", string(debug.Stack())).
			Error("panic recovered")

		if c.Writer.Written() {
			c.Abort()
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.HttpRespBase{
			RequestId: response.RequestId(c),
			Code:      int(t7Error.Unknown.Code),
			Message:   t7Error.Unknown.Message,
		})
	}()
	c.Next()
}

func isBrokenPipe(err error) bool {
	var ne *net.OpError
	if !errors.As(err, &ne) {
		return false
	}
	var se *os.SyscallError
	if !errors.As(ne, &se) {
		return false
	}
	msg := strings.ToLower(se.Error())
	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
}
