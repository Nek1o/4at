package main

import (
	"github.com/gin-contrib/cors"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func (s *ChatServer) configureRoutes() {
	api := s.engine.Group("/api")

	v1 := api.Group("/v1")
	{
		v1.GET("/ping", s.Ping)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	users := v1.Group("/users")
	{
		users.GET("/check/:username", s.UserExists)
	}

	rooms := v1.Group("/rooms")
	rooms.Use(s.Authorization)
	{
		rooms.POST("/:name/", s.AddRoom)
		rooms.DELETE("/:name/", s.RemoveRoom)
		rooms.POST("/join/:name/", s.JoinRoom)
		rooms.POST("/leave/:name/", s.LeaveRoom)
		rooms.GET("/:name", s.GetRoom)
		rooms.GET("/by-owner", s.GetUserRooms)
	}
}

func (s *ChatServer) configureCORS() {
	s.engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"X-App-Token", "X-User-Name", "Origin", ""},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "HEAD", "OPTIONS"},
	}))
}
