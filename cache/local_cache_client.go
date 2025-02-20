package cache

import (
	"context"
	"errors"

	"microblogging/internal/model"
	"microblogging/internal/utils"
)

const maxTweetsSaved = 50

type LocalCacheClient struct {
	cache Cache
}

func NewLocalCacheClient(cache Cache) *LocalCacheClient {
	return &LocalCacheClient{
		cache: cache,
	}
}

func (lc *LocalCacheClient) GetTweetsInCache(ctx context.Context, key string) ([]model.Tweet, error) {
	result, err := lc.cache.Get(key)
	if err != nil && !errors.Is(err, errorKeyNotFound) {
		return nil, err
	}

	if result == nil {
		return []model.Tweet{}, nil
	}

	var savedTweets = make([]model.Tweet, 0)

	err = DecodeCacheValue[[]model.Tweet](result, &savedTweets)
	if err != nil {
		return nil, err
	}

	return savedTweets, nil
}

func (lc *LocalCacheClient) SaveTweetInCache(ctx context.Context, key string, tweet model.Tweet) error {
	tweets, err := lc.GetTweetsInCache(ctx, key)
	if err != nil && !errors.Is(err, errorKeyNotFound) {
		return err
	}

	newTweets := make([]model.Tweet, 0)
	newTweets = append(newTweets, tweet)
	sortedTweets := utils.MergeAndSort[model.Tweet](
		tweets,
		newTweets,
		func(a, b model.Tweet) bool {
			return a.CreatedAt.After(b.CreatedAt)
		})

	if len(sortedTweets) >= maxTweetsSaved {
		sortedTweets = sortedTweets[:maxTweetsSaved]
	}

	resultInBytes, err := EncodeCacheValue[[]model.Tweet](sortedTweets)
	if err != nil {
		return err
	}

	return lc.cache.Set(key, resultInBytes)
}

func (lc *LocalCacheClient) SaveTweetsInCache(ctx context.Context, key string, newTweets []model.Tweet) error {
	tweets, err := lc.GetTweetsInCache(ctx, key)
	if err != nil && !errors.Is(err, errorKeyNotFound) {
		return err
	}

	tweets = append(tweets, newTweets...)

	sortedTweets := utils.MergeAndSort[model.Tweet](
		tweets,
		newTweets,
		func(a, b model.Tweet) bool {
			return a.CreatedAt.After(b.CreatedAt)
		})

	if len(sortedTweets) >= maxTweetsSaved {
		sortedTweets = sortedTweets[:maxTweetsSaved]
	}

	resultInBytes, err := EncodeCacheValue[[]model.Tweet](tweets)
	if err != nil {
		return err
	}

	return lc.cache.Set(key, resultInBytes)
}

func (lc *LocalCacheClient) GetPopularUsersForUserID(ctx context.Context, key string) ([]model.PopularUser, error) {
	result, err := lc.cache.Get(key)
	if err != nil && !errors.Is(err, errorKeyNotFound) {
		return nil, err
	}

	if result == nil {
		return []model.PopularUser{}, nil
	}

	var popularUsers = make([]model.PopularUser, 0)

	err = DecodeCacheValue[[]model.PopularUser](result, &popularUsers)
	if err != nil {
		return nil, err
	}

	return popularUsers, nil
}

func (lc *LocalCacheClient) SavePopularUserForUserID(ctx context.Context, key string, popularUser model.PopularUser) error {
	popularUsers, err := lc.GetPopularUsersForUserID(ctx, key)
	if err != nil && !errors.Is(err, errorKeyNotFound) {
		return err
	}

	popularUsers = append(popularUsers, popularUser)

	resultInBytes, err := EncodeCacheValue[[]model.PopularUser](popularUsers)
	if err != nil {
		return err
	}

	return lc.cache.Set(key, resultInBytes)
}
