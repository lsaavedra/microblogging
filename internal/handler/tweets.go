package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"microblogging/internal/model"
	"microblogging/internal/utils"
)

type (
	TweetsService interface {
		CreateUserTweet(ctx context.Context, tweetRequest model.CreateTweetRequest) (model.CreateTweetResponse, error)
	}

	TimelineSrv interface {
		GetFullUserTimeline(ctx context.Context, userID string) (utils.PaginatedResponse[model.Tweet], error)
	}
)

type TweetsHandler struct {
	service         TweetsService
	timelineService TimelineSrv
}

func NewTweetsHandler(service TweetsService, timelineService TimelineSrv) *TweetsHandler {
	return &TweetsHandler{
		service:         service,
		timelineService: timelineService,
	}
}

func (h *TweetsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/tweets", utils.RouteWithStatus(h.CreateTweet))
	rg.GET("/tweets/timelines", utils.RouteWithStatus(h.GetTweetsTimeline))
}

func (h *TweetsHandler) CreateTweet(c *gin.Context) (int, error) {
	var payload model.CreateTweetRequest
	err := c.BindJSON(&payload)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	if err := payload.Check(); err != nil {
		return http.StatusBadRequest, err
	}

	tweet, err := h.service.CreateUserTweet(c, payload)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to create user tweet")
	}

	c.JSON(http.StatusOK, tweet)

	return http.StatusOK, nil
}

func (h *TweetsHandler) GetTweetsTimeline(c *gin.Context) (int, error) {
	userID := c.Query("user_id")
	if userID == "" {
		return http.StatusBadRequest, errors.New("empty user_id")
	}

	result, err := h.timelineService.GetFullUserTimeline(c, userID)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to get user timeline")
	}

	c.JSON(http.StatusOK, result)

	return http.StatusOK, nil
}
