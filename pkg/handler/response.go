package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type error struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)

	if message == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
		message = "user with this username already exists"
	}
	c.AbortWithStatusJSON(statusCode, error{message})
}
