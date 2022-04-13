package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *ChatServer) Authorization(c *gin.Context) {
	// appHeader := c.Request.Header.Get("x-app-token")
	// TODO add token check
	// if appHeader == TOKEN_CONST {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	usernameHeader := c.Request.Header.Get("x-user-name")
	if usernameHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := s.db.GetUser(c.Request.Context(), usernameHeader)
	if err != nil {
		// if there is no user with such username, then add one
		if errors.Is(err, mongo.ErrNoDocuments) {
			user := &User{Username: usernameHeader}
			s.db.AddUser(c.Request.Context(), user)
			c.Set("user", user)
			c.Next()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}

func (s *ChatServer) GetUser(c *gin.Context) *User {
	ctxUser, ok := c.Get("user")
	if !ok {
		return nil
	}

	user, ok := ctxUser.(*User)
	if !ok {
		return nil
	}

	return user
}
