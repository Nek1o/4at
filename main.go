package main

import (
	"context"
	"log"

	_ "github.com/Nek1o/4at/docs"
	"github.com/gin-gonic/gin"
)

// @title 4at Swagger API
// @version 1.0
// @description This is a 4at API. Contact @Nekiio at telegram for support
// @termsOfService http://swagger.io/terms/

// @contact.name Nekiio
// @contact.email nikita.volchenkov1@gmail.com

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey  Username
// @in header
// @name X-User-Name

// @securityDefinitions.apikey  API token
// @in header
// @name X-App-Name
func main() {
	ctx := context.Background()

	db, err := NewMongoDB(ctx, "localhost", 27017)
	if err != nil {
		log.Fatalf("could not init a db connection: %v", err)
	}

	// gin.SetMode(gin.ReleaseMode)
	app := NewChatServer(gin.Default(), db)

	app.Run()
}

// TODO add x-app-token header
