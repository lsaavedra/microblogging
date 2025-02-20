package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"microblogging/cache"
	"microblogging/internal/model"
	"microblogging/internal/utils"
)

type TimelineService struct {
	logger      *log.Logger
	cacheClient *cache.LocalCacheClient
}

func NewTimelineService(logger *log.Logger, cacheClient *cache.LocalCacheClient) *TimelineService {
	return &TimelineService{
		logger:      logger,
		cacheClient: cacheClient,
	}
}

func (t *TimelineService) GetUserTimeline(ctx context.Context, userID string) (utils.PaginatedResponse[model.Tweet], error) {
	tweets, err := t.cacheClient.GetTweetsInCache(ctx, userID)
	if err != nil {
		return utils.BuildResponse[model.Tweet](tweets), err
	}

	log.Infof("Get user timeline for user: %v - total tweets: %v", userID, len(tweets))

	return utils.BuildResponse[model.Tweet](tweets), nil
}

func (t *TimelineService) GetFullUserTimeline(ctx context.Context, userID string) (utils.PaginatedResponse[model.Tweet], error) {
	tweets, err := t.getUserTimelineTweets(ctx, userID)
	if err != nil {
		return utils.BuildResponse[model.Tweet](tweets), err
	}

	log.Infof("Get user timeline for user: %v - total tweets: %v", userID, len(tweets))

	return utils.BuildResponse[model.Tweet](tweets), nil
}

func (t *TimelineService) UpdateUserTimeline(ctx context.Context, userID string, updatedTweets []model.Tweet) error {
	return t.cacheClient.SaveTweetsInCache(ctx, userID, updatedTweets)
}

func (t *TimelineService) getUserTimelineTweets(ctx context.Context, userID string) ([]model.Tweet, error) {
	tweets, err := t.cacheClient.GetTweetsInCache(ctx, userID)
	if err != nil {
		return nil, err
	}

	popularTweets := make([]model.Tweet, 0)
	key := fmt.Sprintf(model.KeyForPopularUsers, userID)
	popularUsers, err := t.cacheClient.GetPopularUsersForUserID(ctx, key)
	if len(popularUsers) != 0 {
		log.Infof("agreggating global tweets from:%v users", len(popularUsers))
		for _, popularUser := range popularUsers {
			popularTweetsForUserKey := fmt.Sprintf(model.KeyForPopularTweets, popularUser.UserID)
			log.Infof("getting global tweets for user key:%v", popularTweetsForUserKey)
			popularTweetsFromUser, err := t.cacheClient.GetTweetsInCache(ctx, popularTweetsForUserKey)
			if err != nil {
				continue
			}
			popularTweets = append(popularTweets, popularTweetsFromUser...)
		}
	}

	tweets = append(tweets, popularTweets...)

	tweets = utils.MergeAndSort(tweets, popularTweets, func(a, b model.Tweet) bool {
		return a.CreatedAt.After(b.CreatedAt)
	})

	if len(tweets) > 50 {
		return tweets[:49], nil
	}

	return tweets, nil
}
