package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func (m *Controller) RecoverMiddleware(c *gin.Context) {
	// TODO: refine
	defer func() {
		if err := recover(); err != nil {
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}

			var brokenPipe bool
			var ne *net.OpError
			if errors.As(err.(error), &ne) {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			if brokenPipe {
				_ = c.Error(err.(error))
				c.Abort()
			} else {
				m.log.WithService("gin").
					With("result", err).
					With("header", headers).
					Panic("panic recovered")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	c.Next()
}
