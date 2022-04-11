package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *ChatServer) Authentication(c *gin.Context) {
	header := strings.Split(c.Request.Header.Get("authorization"), " ")
	if len(header) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	if header[0] != "Bearer" && header[1] == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header value, shoud be \"Bearer <token>\""})
		c.Abort()
		return
	}

	token := header[1]

	user, err := s.db.GetUser(c.Request.Context(), "", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
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
