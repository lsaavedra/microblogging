package service

import (
	"context"

	log "github.com/sirupsen/logrus"

	"microblogging/internal/eventstub"
	"microblogging/internal/model"
)

type TweetsService struct {
	logger    *log.Logger
	repo      TweetsRepo
	publisher eventstub.EventProcessor
}

type TweetsRepo interface {
	CreateTweet(ctx context.Context, tweet model.Tweet) (model.Tweet, error)
	GetByFilters(
		filters map[string]interface{},
		limit, offset int,
		orderBy string,
	) ([]model.Tweet, error)
}

func NewTweetsService(logger *log.Logger, repo TweetsRepo, publisher eventstub.EventProcessor) *TweetsService {
	return &TweetsService{
		logger:    logger,
		repo:      repo,
		publisher: publisher,
	}
}

func (t *TweetsService) CreateUserTweet(ctx context.Context, tweetRequest model.CreateTweetRequest) (model.CreateTweetResponse, error) {
	tweet := model.Tweet{
		CreatedBy:       tweetRequest.UserID,
		Message:         tweetRequest.Message,
		MediaContentUrl: tweetRequest.MediaContentUrl,
	}

	log.Info("Creating tweet from user:", tweetRequest.UserID)
	tweet, err := t.repo.CreateTweet(ctx, tweet)
	if err != nil {
		return model.CreateTweetResponse{}, err
	}

	log.Info("Publishing event for tweet created")
	t.publisher.Publish(model.TweetCreated, model.Event{Data: tweet})

	return model.CreateTweetResponse{ID: tweet.ID}, nil
}

func (t *TweetsService) GetRecentUserTweets(ctx context.Context, userID string) ([]model.Tweet, error) {
	filters := make(map[string]interface{})
	filters["created_by"] = userID

	//Will only get last 10 tweets from user
	recentTweets, err := t.repo.GetByFilters(filters, 10, 0, "created_at DESC")
	if err != nil {
		return nil, err
	}

	return recentTweets, nil
}
