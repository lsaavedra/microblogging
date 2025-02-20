package service

import (
	"context"

	"github.com/google/uuid"

	"microblogging/internal/eventstub"
	"microblogging/internal/model"
)

type FollowsService struct {
	repo      FollowersRepo
	publisher eventstub.EventProcessor
}

type FollowersRepo interface {
	CreateFollowing(ctx context.Context, follow model.Follow) error
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
}

func NewFollowsService(repo FollowersRepo, publisher eventstub.EventProcessor) *FollowsService {
	return &FollowsService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *FollowsService) FollowUser(ctx context.Context, userID, targetUserID uuid.UUID) error {
	follow := model.Follow{
		FollowedID: targetUserID,
		FollowerID: userID,
	}

	err := s.repo.CreateFollowing(ctx, follow)
	if err != nil {
		return err
	}

	s.publisher.Publish(model.UserFollowed, model.Event{Data: follow})

	return nil
}

func (s *FollowsService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	result, err := s.repo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *FollowsService) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int, error) {
	result, err := s.GetFollowers(ctx, userID)
	if err != nil {
		return 0, err
	}

	return len(result), nil
}
