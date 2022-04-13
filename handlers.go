package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *ChatServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (s *ChatServer) AddUser(c *gin.Context) {
	var addUser AddUser
	if err := c.ShouldBindJSON(&addUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := User{Username: addUser.Username}
	if err := s.db.AddUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *ChatServer) UserExists(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	exists, err := s.db.UserExists(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (s *ChatServer) AddRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room name"})
		return
	}

	if room, err := s.db.AddRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, room)
	}
}

func (s *ChatServer) RemoveRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if err := s.db.RemoveRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *ChatServer) JoinRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if room, err := s.db.JoinRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, room)
	}
}

func (s *ChatServer) LeaveRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if err := s.db.LeaveRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *ChatServer) GetRoom(c *gin.Context) {
	roomName := c.Param("name")

	var room *Room
	room, err := s.db.GetRoom(c.Request.Context(), roomName)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (s *ChatServer) GetUserRooms(c *gin.Context) {
	user := s.GetUser(c)

	if rooms, err := s.db.GetUserRooms(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, rooms)
	}

}
