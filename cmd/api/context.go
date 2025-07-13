package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rizky-ardiansah/event-api/internal/database"
)

func (app *application) GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &database.User{}
	}
	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}
	return user
}
