package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/Template7/backend/docs"
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
	app := InitializeApp()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Service.Port),
		Handler: app.SetupRoutes(),
	}

	go func() {
		app.Log.With("port", app.config.Service.Port).Info("server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Log.WithError(err).Panic(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Log.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.Log.WithError(err).Panic("fail to shutdown server")
	}

	app.Log.Info("server exited properly")
}
