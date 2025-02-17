package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microblogging/internal/utils"
)

type TweetsService interface {
	//TODO: define methods
}

type TweetsHandler struct {
	service TweetsService
}

func NewTweetsHandler(service UserService) *TweetsHandler {
	return &TweetsHandler{
		service: service}
}

func (h *TweetsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/tweet", utils.RouteWithStatus(h.CreateTweet))
	rg.PATCH("/tweet/:id", utils.RouteWithStatus(h.UpdateTweetByID))
	rg.GET("/tweets", utils.RouteWithStatus(h.GetTweetsTimeline))
}

func (h *TweetsHandler) CreateTweet(c *gin.Context) (int, error) {
	//TODO: handle request filters and called service
	//Check that tweet is less than 280 characters

	c.JSON(http.StatusOK, map[string]string{})

	return http.StatusOK, nil
}

func (h *TweetsHandler) UpdateTweetByID(c *gin.Context) (int, error) {
	//TODO: handle request filters and called service

	c.JSON(http.StatusOK, map[string]string{})

	return http.StatusOK, nil
}

func (h *TweetsHandler) GetTweetsTimeline(c *gin.Context) (int, error) {
	//TODO: handle request filters and called service

	c.JSON(http.StatusOK, map[string]string{})

	return http.StatusOK, nil
}
