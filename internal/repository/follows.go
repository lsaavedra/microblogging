package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"

	"microblogging/db"
	"microblogging/internal/model"
)

const errorUniqueViolationPgCode = "23505"

var ErrFollowRelationAlreadyExists = errors.New("follow relation already exists")

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

func (repo *FollowsRepository) CreateFollowing(ctx context.Context, follow model.Follow) error {
	follow.FollowedAt = time.Now()

	err := repo.db.WithContext(ctx).Create(&follow).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == errorUniqueViolationPgCode {
			return ErrFollowRelationAlreadyExists
		}
	}

	return err
}

func (repo *FollowsRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var following []model.Follow

	err := repo.db.WithContext(ctx).Where("followed_id = ?", userID).Find(&following).Error

	return following, err
}
