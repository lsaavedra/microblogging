package repository

import (
	"context"

	"github.com/google/uuid"

	"microblogging/db"
	"microblogging/internal/model"
)

type TweetsRepository struct {
	db *db.Database
}

func NewTweetRepository(db *db.Database) *TweetsRepository {
	return &TweetsRepository{
		db: db,
	}
}

func (repo *TweetsRepository) CreateTweet(ctx context.Context, tweet model.Tweet) (uuid.UUID, error) {
	err := repo.db.WithContext(ctx).Create(&tweet).Error

	return tweet.ID, err
}
