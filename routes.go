package main

func (s *ChatServer) configureRoutes() {
	api := s.engine.Group("/api")

	v1 := api.Group("/v1")
	{
		v1.GET("/ping", s.Ping)
	}

	users := v1.Group("/users")
	{
		users.POST("/add", s.AddUser)
		users.POST("/check", s.UserExists)
	}

	rooms := v1.Group("/rooms")
	rooms.Use(s.Authorization)
	{
		rooms.POST("/add", s.AddRoom)
		rooms.POST("/remove", s.RemoveRoom)
		rooms.POST("/join", s.JoinRoom)
		rooms.POST("/leave", s.LeaveRoom)
		rooms.POST("/get", s.GetRoom)
		rooms.GET("/get-by-owner", s.GetUserRooms)
	}
}
