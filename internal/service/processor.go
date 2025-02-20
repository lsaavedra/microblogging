package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"microblogging/cache"
	"microblogging/internal/model"
	"microblogging/internal/utils"
)

const (
	maxFollowersProcess = 10000 // Threshold to know how many followers the user could have to update follower timelines
)

type (
	ProcessorService struct {
		logger           *log.Logger
		cacheClient      *cache.LocalCacheClient
		followersService FollowersSrv
		tweetsService    TweetsSrv
		timelineService  TimelineSrv
	}

	FollowersSrv interface {
		GetFollowersCount(ctx context.Context, userID uuid.UUID) (int, error)
		GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
	}

	TweetsSrv interface {
		GetRecentUserTweets(ctx context.Context, userID string) ([]model.Tweet, error)
	}

	TimelineSrv interface {
		GetUserTimeline(ctx context.Context, userID string) (utils.PaginatedResponse[model.Tweet], error)
		UpdateUserTimeline(ctx context.Context, userID string, updatedTweets []model.Tweet) error
	}
)

func NewProcessorService(
	logger *log.Logger,
	followerService FollowersSrv,
	tweetsService TweetsSrv,
	timelineService TimelineSrv,
	cacheClient *cache.LocalCacheClient,
) *ProcessorService {
	srv := &ProcessorService{
		logger:           logger,
		followersService: followerService,
		tweetsService:    tweetsService,
		timelineService:  timelineService,
		cacheClient:      cacheClient,
	}

	return srv
}

func (p *ProcessorService) ProcessIncomingTweet(ctx context.Context, tweet model.Tweet) error {
	followers, err := p.followersService.GetFollowers(ctx, tweet.CreatedBy)
	if err != nil {
		return err
	}

	log.Infof("[incoming_tweet] user_id: %v with total followers: %v", tweet.CreatedBy.String(), len(followers))
	if len(followers) == 0 {
		return nil
	}

	if len(followers) < maxFollowersProcess {
		for _, followRelation := range followers {
			log.Infof("[incoming_tweet] updating timeline cache for user_id: %v", followRelation.FollowerID.String())
			err := p.cacheClient.SaveTweetInCache(ctx, followRelation.FollowerID.String(), tweet)
			if err != nil {
				log.Errorf("error updating timeline for user: %v", followRelation.FollowerID.String())
			}
		}
	} else {
		keyPopularForUser := fmt.Sprintf(model.KeyForPopularTweets, tweet.CreatedBy)
		log.Infof(
			"[incoming_tweet] popular user: %v - saving tweets in popular cache: %v",
			tweet.CreatedBy.String(), keyPopularForUser,
		)
		err = p.cacheClient.SaveTweetInCache(ctx, keyPopularForUser, tweet)
		if err != nil {
			log.Errorf("error saving tweet for global user in cache: %v", tweet.CreatedBy)
		}
	}

	return nil
}

func (p *ProcessorService) ProcessNewFollower(ctx context.Context, followRelation model.Follow) error {
	var tweetsFromFollowed []model.Tweet

	followedID := followRelation.FollowedID
	followerID := followRelation.FollowerID
	total, err := p.followersService.GetFollowersCount(ctx, followedID)
	if err != nil {
		return err
	}
	log.Infof("[new_follower] followed_id: %v with total followers: %v", followedID, total)

	if total == 0 {
		return nil
	}

	if total < maxFollowersProcess {
		tweetsFromFollowed, err = p.tweetsService.GetRecentUserTweets(ctx, followedID.String())
		if err != nil {
			return err
		}
		log.Infof("[new_follower] getting recent tweets: %v from followed_id: %v", len(tweetsFromFollowed), followedID)
		if len(tweetsFromFollowed) == 0 {
			return nil
		}
		log.Infof("[new_follower] adding to timeline: %v tweets for follower_id: %v", len(tweetsFromFollowed), followerID)
		err = p.timelineService.UpdateUserTimeline(ctx, followerID.String(), tweetsFromFollowed)
		if err != nil {
			return err
		}
	} else {
		key := fmt.Sprintf(model.KeyForPopularUsers, followerID.String())
		log.Infof(
			"[new_follower] popular user - updating only popular list key: %v from follower_id: %v",
			key, followerID,
		)
		err = p.cacheClient.SavePopularUserForUserID(ctx, key, model.PopularUser{UserID: followedID.String()})
		if err != nil {
			return err
		}
	}

	return nil

}

func (p *ProcessorService) FollowersProcessor(event model.Event) {
	log.Infof("[topic %s] reading event ...", model.UserFollowed)
	follow := event.Data.(model.Follow)
	err := p.ProcessNewFollower(context.TODO(), follow)
	if err != nil {
		log.Error("error processing new follower event")
	}

}

func (p *ProcessorService) TweetsProcessor(event model.Event) {
	log.Infof("[topic %s] reading event ...", model.TweetCreated)
	tweet := event.Data.(model.Tweet)
	err := p.ProcessIncomingTweet(context.TODO(), tweet)
	if err != nil {
		log.Error("error processing new tweet created")
	}

}
