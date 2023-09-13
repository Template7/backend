package middleware

import (
	"errors"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func RecoverMiddleware(c *gin.Context) {
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
			if ne, ok := err.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			if brokenPipe {
				c.Error(err.(error))
				c.Abort()
			} else {
				logger.New().WithService("gin").
					With("result", err).
					With("header", headers).
					Panic("panic recovered")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	c.Next()
}
