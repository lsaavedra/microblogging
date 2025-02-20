package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	KeyForPopularUsers = "follower_populars|%s"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type UserFollowingRequestBody struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

func (u *UserFollowingRequestBody) Check() error {
	if u.UserID == u.TargetUserID {
		return errors.New("both users_ids are equal")
	}

	return nil
}

type PopularUser struct {
	UserID string
}
