package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

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
// TODO add x-user-name header

// TODO rework auth process
// TODO rework routes
// TODO add swagger