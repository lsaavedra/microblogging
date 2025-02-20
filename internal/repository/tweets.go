package repository

import (
	"context"

	"gorm.io/gorm"

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

func (repo *TweetsRepository) CreateTweet(ctx context.Context, tweet model.Tweet) (model.Tweet, error) {
	err := repo.db.WithContext(ctx).Create(&tweet).Error

	return tweet, err
}

func (repo *TweetsRepository) GetByFilters(
	filters map[string]interface{},
	limit, offset int,
	orderBy string,
) ([]model.Tweet, error) {
	var results []model.Tweet

	query := repo.buildQuery(filters, limit, offset, orderBy)
	result := query.Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

func (repo *TweetsRepository) buildQuery(filters map[string]interface{}, limit, offset int, orderBy string) *gorm.DB {
	query := repo.db.DB

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	return query.Limit(limit).Offset(offset)
}
