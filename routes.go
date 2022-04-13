package main

func (s *ChatServer) configureRoutes() {
	api := s.engine.Group("/api")

	v1 := api.Group("/v1")
	{
		v1.GET("/ping", s.Ping)
	}

	users := v1.Group("/users")
	{
		users.POST("/check/:username/", s.UserExists)
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
