package service

import (
	"context"

	"github.com/google/uuid"

	"microblogging/internal/model"
)

type UsersService struct {
	Repo UsersRepo
}

func (u UsersService) GetFollowingUsers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	//TODO implement me
	panic("implement me")
}

func (u UsersService) FollowUser(ctx context.Context, userID, targetUserID uuid.UUID) error {
	follow := model.Follow{
		FollowerID:  userID,
		FollowingID: targetUserID,
	}
	return u.Repo.SaveFollowAndFollowing(ctx, follow)
}

type UsersRepo interface {
	SaveFollowAndFollowing(ctx context.Context, follow model.Follow) error
}

func NewUsersService(repo UsersRepo) *UsersService {
	return &UsersService{
		Repo: repo,
	}
}
