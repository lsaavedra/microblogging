package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AbortError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func NewCustomerUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := verifyUserAuth(ctx)

		if err != nil {
			abortWithError(ctx, http.StatusUnauthorized, err)
			return
		}

		ctx.Next()
	}
}

func verifyUserAuth(ctx *gin.Context) error {
	// TODO we should validate that x-user-id comes from headers or in body
	/*
		userId := ctx.Request.Header.Get("x-user-id")

		if userId == "" {
			return errors.New("user_id is required")
		}
	*/

	//TODO add logic to resolve that userid comes from header, or body o param

	return nil
}

func abortWithError(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, AbortError{
		StatusCode: statusCode,
		Message:    err.Error(),
	})
}
