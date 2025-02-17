package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microblogging/db"
	"microblogging/internal/handler/middleware"
	"microblogging/internal/repository"
	"microblogging/internal/service"
)

func RouterWithHandlers(database db.Database) *gin.Engine {
	router := gin.Default()

	router.GET("/", healthHandler)

	globalMiddleware := middleware.NewCustomerUserMiddleware()

	apiV1Group := router.Group("/api/v1", globalMiddleware)

	userRepository := repository.NewFollowsRepository(database)
	userService := service.NewUsersService(userRepository)
	userHandler := NewUserHandler(userService)
	userHandler.RegisterRoutes(apiV1Group)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
