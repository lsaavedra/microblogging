package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID  uuid.UUID `gorm:"primaryKey;type:uuid"`
	FollowingID uuid.UUID `gorm:"primaryKey"`
	FollowedAt  time.Time
}
