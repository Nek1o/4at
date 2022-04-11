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

	app := NewChatServer(gin.Default(), db)

	app.Run()
}
