package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
