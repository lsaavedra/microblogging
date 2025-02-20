package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AbortError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func NewCustomerUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := verifyUserFromRequest(ctx)
		if err == nil {
			ctx.Next()
			return
		}

		if err != nil {
			abortWithError(ctx, http.StatusUnauthorized, err)
			return
		}

		ctx.Next()
	}
}

func verifyUserFromRequest(ctx *gin.Context) error {
	userIdFromParams := ctx.Query("user_id")
	if userIdFromParams == "" {
		userIdFromHeaders := ctx.Request.Header.Get("x-user-id")
		if userIdFromHeaders == "" {
			return errors.New("user_id is required by param or header")
		}
	}

	return nil
}

func abortWithError(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, AbortError{
		StatusCode: statusCode,
		Message:    err.Error(),
	})
}
