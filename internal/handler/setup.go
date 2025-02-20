package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"microblogging/cache"
	"microblogging/db"
	"microblogging/internal/eventstub"
	"microblogging/internal/handler/middleware"
	"microblogging/internal/model"
	"microblogging/internal/repository"
	"microblogging/internal/service"
)

func RouterWithHandlers(database *db.Database, logger *log.Logger, pubSub eventstub.EventProcessor) *gin.Engine {
	router := gin.Default()

	router.GET("/", healthHandler)

	globalMiddleware := middleware.NewCustomerUserMiddleware()

	apiV1Group := router.Group("/api/v1", globalMiddleware)

	followersRepository := repository.NewFollowsRepository(database)
	followersService := service.NewFollowsService(followersRepository, pubSub)
	followersHandler := NewFollowersHandler(followersService)
	followersHandler.RegisterRoutes(apiV1Group)

	localCache := cache.NewLocalCache()
	cacheClient := cache.NewLocalCacheClient(localCache)

	tweetsRepository := repository.NewTweetRepository(database)
	tweetsService := service.NewTweetsService(logger, tweetsRepository, pubSub)
	timelineService := service.NewTimelineService(logger, cacheClient)
	tweetsHandler := NewTweetsHandler(tweetsService, timelineService)
	tweetsHandler.RegisterRoutes(apiV1Group)

	processorService := service.NewProcessorService(logger, followersService, tweetsService, timelineService, cacheClient)

	logger.Info("Subscribing to events for follower users and tweets creation ...")
	pubSub.Subscribe(model.UserFollowed, processorService.FollowersProcessor)
	pubSub.Subscribe(model.TweetCreated, processorService.TweetsProcessor)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
