package main

import (
	"context"
	"errors"
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

	r := gin.New()
	route.Setup(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.New().Service.Port),
		Handler: r,
	}

	go func() {
		log.With("port", config.New().Service.Port).Info("server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Panic(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Panic("fail to shutdown server")
	}

	log.Info("server exited properly")
}
