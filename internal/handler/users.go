package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"microblogging/internal/model"
	"microblogging/internal/utils"
)

type UserService interface {
	FollowUser(ctx context.Context, userID, targetUserID uuid.UUID) error
	GetFollowingUsers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users/follow", utils.RouteWithStatus(h.AddUserFollowing))
	rg.GET("/users/follow", utils.RouteWithStatus(h.GetUsersFollowing))
}

func (h *UserHandler) AddUserFollowing(c *gin.Context) (int, error) {
	var payload model.UserFollowingRequestBody
	err := c.BindJSON(&payload)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	err = h.service.FollowUser(c, payload.UserID, payload.TargetUserID)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to follow new user")
	}

	c.JSON(http.StatusOK, map[string]string{})

	return http.StatusOK, nil
}

func (h *UserHandler) GetUsersFollowing(c *gin.Context) (int, error) {
	result, err := h.service.GetFollowingUsers(c, uuid.New())
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to get following users")
	}

	c.JSON(http.StatusOK, result)

	return http.StatusOK, nil
}
