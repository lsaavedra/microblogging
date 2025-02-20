package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"microblogging/internal/model"
)

func TestLocalCacheClient_GetTweetsInCacheWhenNoValuePresent(t *testing.T) {
	c := NewLocalCache()
	cacheClient := NewLocalCacheClient(c)

	result, err := cacheClient.GetTweetsInCache(context.TODO(), "12345")
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestLocalCacheClient_GetTweetsInCacheWithValuesPresent(t *testing.T) {
	c := NewLocalCache()
	cacheClient := NewLocalCacheClient(c)
	key := "user_2240"
	userID := uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")

	result, err := cacheClient.GetTweetsInCache(context.TODO(), key)
	assert.NoError(t, err)
	assert.Empty(t, result)

	_ = cacheClient.SaveTweetInCache(context.TODO(), key, model.Tweet{
		CreatedBy: userID,
		Message:   "My first tweet",
		CreatedAt: time.Now().Add(-60 * time.Second),
	})

	result, err = cacheClient.GetTweetsInCache(context.TODO(), key)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
