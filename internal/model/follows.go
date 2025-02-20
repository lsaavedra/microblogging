package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowedID uuid.UUID `gorm:"primaryKey;type:uuid"`
	FollowerID uuid.UUID `gorm:"primaryKey;type:uuid"`
	FollowedAt time.Time
}
