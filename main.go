package main

import (
	"context"
	"fmt"
	_ "github.com/Template7/backend/docs"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/route"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Backend API
// @version 1.0
// @description API Documentation

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// schemes http
func main() {
	log := logger.New().WithService("main")

	r := gin.Default()

	gin.SetMode(config.New().Gin.Mode)
	route.Setup(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.New().Gin.ListenPort),
		Handler: r,
	}

	go func() {
		log.With("port", config.New().Gin.ListenPort).Info("server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Panic(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.New().Gin.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Panic("fail to shutdown server")
	}

	log.Info("server exited properly")
}
