package repository

import (
	"context"

	"microblogging/db"
	"microblogging/internal/model"
)

type TimelineRepository struct {
	db *db.Database
}

func NewTimelineRepository(db *db.Database) *TimelineRepository {
	return &TimelineRepository{
		db: db,
	}
}

func (repo *TimelineRepository) AddTweetToUserTimeline(ctx context.Context, userID string, tweet model.Tweet) error {
	//GET redis cache user timeline
	//Update or append to timeline
	return nil
}
