package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ping godoc
// @Summary Ping the server.
// @Description Get the status of server.
// @Tags Ping
// @Accept */*
// @Produce json
// @Success 200
// @Router /ping [get]
func (s *ChatServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// AddUser is deprecated
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

// UserExists godoc
// @Summary Check that the user exists
// @Description Check that the user exists.
// @Tags User
// @Accept */*
// @Produce json
// @Param username path string false "username to check" default(Tolya)
// @Success 200 {object} CheckUser
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/check/{username} [get]
func (s *ChatServer) UserExists(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid username"))
		return
	}

	exists, err := s.db.UserExists(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, CheckUser{exists})
}

// AddRoom godoc
// @Summary Add a new chat room
// @Description Add new a chat room.
// @Tags Room
// @Accept */*
// @Produce json
// @Param name path string false "room name" default(Tolya's room)
// @Success 200 {object} Room
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rooms/{name}/ [post]
// @Security Username
// @Security API token
func (s *ChatServer) AddRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid room name"))
		return
	}

	if room, err := s.db.AddRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, room)
	}
}

// RemoveRoom godoc
// @Summary Remove a room from chat
// @Description Remove a room from chat.
// @Tags Room
// @Accept */*
// @Produce json
// @Param name path string false "room name" default(Tolya's room)
// @Success 200
// @Failure 500 {object} ErrorResponse
// @Router /rooms/{name}/ [delete]
// @Security Username
// @Security API token
func (s *ChatServer) RemoveRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if err := s.db.RemoveRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// JoinRoom godoc
// @Summary Join a chat room
// @Description Join a chat room.
// @Tags Room
// @Accept */*
// @Produce json
// @Param name path string false "room name" default(Tolya's room)
// @Success 200 {object} Room
// @Failure 500 {object} ErrorResponse
// @Router /rooms/join/{name}/ [post]
// @Security Username
// @Security API token
func (s *ChatServer) JoinRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if room, err := s.db.JoinRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, room)
	}
}

// LeaveRoom godoc
// @Summary Leave a chat room
// @Description Leave a chat room.
// @Tags Room
// @Accept */*
// @Produce json
// @Param name path string false "room name" default(Tolya's room)
// @Success 200 {object} Room
// @Failure 500 {object} ErrorResponse
// @Router /rooms/leave/{name}/ [post]
// @Security Username
// @Security API token
func (s *ChatServer) LeaveRoom(c *gin.Context) {
	user := s.GetUser(c)
	roomName := c.Param("name")

	if err := s.db.LeaveRoom(c.Request.Context(), roomName, user); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// GetRoom godoc
// @Summary Get a chat room
// @Description Get a chat room.
// @Tags Room
// @Accept */*
// @Produce json
// @Param name path string false "room name" default(Tolya's room)
// @Success 200 {object} Room
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rooms/{name} [get]
// @Security Username
// @Security API token
func (s *ChatServer) GetRoom(c *gin.Context) {
	roomName := c.Param("name")

	var room *Room
	room, err := s.db.GetRoom(c.Request.Context(), roomName)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, room)
}

// GetUserRooms godoc
// @Summary Get user's chat rooms
// @Description Get user's chat rooms.
// @Tags Room
// @Accept */*
// @Produce json
// @Success 200 {object} []Room
// @Failure 500 {object} ErrorResponse
// @Router /rooms/by-owner [get]
// @Security Username
// @Security API token
func (s *ChatServer) GetUserRooms(c *gin.Context) {
	user := s.GetUser(c)

	if rooms, err := s.db.GetUserRooms(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, rooms)
	}

}
