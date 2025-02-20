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

type FollowsService interface {
	FollowUser(ctx context.Context, userID, targetUserID uuid.UUID) error
}

type FollowersHandler struct {
	service FollowsService
}

func NewFollowersHandler(service FollowsService) *FollowersHandler {
	return &FollowersHandler{
		service: service}
}

func (h *FollowersHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users/follow", utils.RouteWithStatus(h.AddUserFollowing))
}

func (h *FollowersHandler) AddUserFollowing(c *gin.Context) (int, error) {
	var payload model.UserFollowingRequestBody
	err := c.BindJSON(&payload)

	if err != nil {
		return http.StatusBadRequest, nil
	}

	if err = payload.Check(); err != nil {
		return http.StatusBadRequest, err
	}

	err = h.service.FollowUser(c, payload.UserID, payload.TargetUserID)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to follow new user")
	}

	c.JSON(http.StatusOK, map[string]string{})

	return http.StatusOK, nil
}
