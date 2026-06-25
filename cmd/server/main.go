package main

import (
	"context"
	"sync"

	"github.com/youknow2509/temp-go-ddd/internal/initialize"
)

// @title           Swagger API
// @version         1.0
// @description     This is a REST API for system base.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/Youknow2509/
// @contact.email  lytranvinh.work@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Initialize application context and wait group
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start system
	resources, err := initialize.Initialize(ctx, wg)
	if err != nil {
		panic(err)
	}

	// Wait for OS shutdown signal and gracefully shut down services
	initialize.WaitForShutdown(cancel, wg, resources)
}
