package main

import (
	"sync"

	"github.com/youknow2509/temp-go-ddd/internal/global"
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
	// Inintialize global wait group
	global.WaitGroup = &sync.WaitGroup{}

	// Start system
	if err := initialize.Initialize(); err != nil {
		panic(err)
	}

	// Wait for all goroutines to finish
	global.WaitGroup.Wait()
}
