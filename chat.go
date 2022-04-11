package main

import "github.com/gin-gonic/gin"

type ChatServer struct {
	engine *gin.Engine
	db DB
}

func NewChatServer(engine *gin.Engine, db DB) ChatServer {
	s := ChatServer{engine, db}
	s.configureRoutes() 
	return s
}

func (s *ChatServer) Run(addr ...string) error {
	return s.engine.Run(addr...)
}