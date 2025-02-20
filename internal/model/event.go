package model

const (
	UserFollowed = "user_followed"
	TweetCreated = "tweet_created"
)

type Event struct {
	Data interface{}
}
