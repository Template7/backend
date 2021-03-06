package main

import (
	"context"
	"fmt"
	_ "github.com/Template7/backend/docs"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/route"
	"github.com/Template7/backend/internal/pkg/t7Redis"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	log = logger.GetLogger()
)

// @title Backend API
// @version 1.0
// @description API Documentation

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// schemes http
func main() {

	r := gin.Default()

	gin.SetMode(config.New().Gin.Mode)
	route.Setup(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.New().Gin.ListenPort),
		Handler: r,
	}

	go func() {
		log.Debug("server listen on: ", config.New().Gin.ListenPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
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
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.New().Gin.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Info("server exited properly")
}

func init() {
	config.New()
	db.New()
	t7Redis.New()
}
