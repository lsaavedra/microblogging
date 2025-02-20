package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const KeyForPopularTweets = "popular|%s"

type (
	Tweet struct {
		ID              uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid"`
		CreatedBy       uuid.UUID      `json:"created_by"`
		Message         string         `json:"message"`
		MediaContentUrl string         `json:"media_content_url"`
		CreatedAt       time.Time      `json:"created_at"`
		UpdatedAt       time.Time      `json:"updated_at"`
		DeletedAt       gorm.DeletedAt `json:"deleted_at"`
	}

	CreateTweetRequest struct {
		UserID          uuid.UUID `json:"user_id" bind:"required"`
		Message         string    `json:"message" bind:"required"`
		MediaContentUrl string
	}

	CreateTweetResponse struct {
		ID uuid.UUID `json:"id"`
	}

	TweetDTO struct {
		CreatedBy uuid.UUID `json:"created_by"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func (base *Tweet) BeforeCreate(tx *gorm.DB) error {
	var idValue uuid.UUID

	if base.ID == uuid.Nil {
		idValue = uuid.New()
		tx.Statement.SetColumn("ID", idValue)
	}

	return nil
}

func (c *CreateTweetRequest) Check() error {
	const maxMessageLen = 280

	if len(c.Message) > maxMessageLen {
		return errors.New("message is too long")
	}

	return nil
}

func (t Tweet) ToApi() TweetDTO {
	return TweetDTO{
		CreatedBy: t.CreatedBy,
		Message:   t.Message,
		CreatedAt: t.CreatedAt,
	}
}
