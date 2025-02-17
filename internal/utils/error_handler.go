package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HandlerWithStatus func(c *gin.Context) (int, error)

const defaultMessage = "ERROR"

func RouteWithStatus(handleFunc HandlerWithStatus) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, err := handleFunc(c)
		verifyError(c, status, err)
	}
}

func verifyError(c *gin.Context, status int, err error) {
	if err != nil {
		err := errors.Wrap(err, defaultMessage)
		err = c.Error(err)
		c.JSON(status, gin.H{"error": err.Error()})

		return
	}
}
