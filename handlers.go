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
	var getUser GetUser
	if err := c.ShouldBindJSON(&getUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := s.db.UserExists(c.Request.Context(), getUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (s *ChatServer) AddRoom(c *gin.Context) {
	user := s.GetUser(c)

	var addRoom AddRoom
	if err := c.ShouldBindJSON(&addRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if room, err := s.db.AddRoom(c.Request.Context(), addRoom.Name, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, room)
	}
}

func (s *ChatServer) RemoveRoom(c *gin.Context) {
	user := s.GetUser(c)

	var removeRoom RemoveRoom
	if err := c.ShouldBindJSON(&removeRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.RemoveRoom(c.Request.Context(), removeRoom.Name, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *ChatServer) JoinRoom(c *gin.Context) {
	user := s.GetUser(c)

	var joinRoom JoinRoom
	if err := c.ShouldBindJSON(&joinRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if room, err := s.db.JoinRoom(c.Request.Context(), joinRoom.Name, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, room)
	}
}

func (s *ChatServer) LeaveRoom(c *gin.Context) {
	user := s.GetUser(c)

	var leaveRoom LeaveRoom
	if err := c.ShouldBindJSON(&leaveRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.LeaveRoom(c.Request.Context(), leaveRoom.Name, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *ChatServer) GetRoom(c *gin.Context) {
	var getRoom GetRoom
	if err := c.ShouldBindJSON(&getRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var room *Room
	room, err := s.db.GetRoom(c.Request.Context(), getRoom.Name)
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
