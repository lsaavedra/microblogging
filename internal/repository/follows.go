package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"microblogging/db"
	"microblogging/internal/model"
)

type (
	FollowsRepository struct {
		db *db.Database
	}
)

func NewFollowsRepository(db *db.Database) *FollowsRepository {
	return &FollowsRepository{
		db: db,
	}
}

func (repo *FollowsRepository) SaveFollowAndFollowing(ctx context.Context, follow model.Follow) error {
	//TODO resolve how to connect to db to add new relation
	fmt.Printf("User_id: %v following a new user_id: %v", follow.FollowerID, follow.FollowingID)

	follow.FollowedAt = time.Now()
	err := repo.db.WithContext(ctx).Create(&follow).Error

	return err
}

func (repo *FollowsRepository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var following []uuid.UUID

	err := repo.db.WithContext(ctx).Where("follower_id = ?", userID).Find(&following).Error

	return following, err
}
