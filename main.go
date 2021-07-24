package main

import (
	_ "backend/docs"
	"backend/internal/pkg/config"
	"backend/internal/pkg/route"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @title Backend API
// @version 1.0
// @description API Documentation

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// schemes http
func main() {

	r := gin.Default()
	gin.SetMode(config.New().Mode)
	route.Setup(r)

	// start http server and listen on default port 8080
	_ = r.Run(fmt.Sprintf(":%d", config.New().Port))
}
