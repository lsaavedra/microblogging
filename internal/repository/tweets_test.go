package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"microblogging/db"
	"microblogging/internal/model"
)

var (
	createdBy = uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
	tweet     = model.Tweet{
		ID:              uuid.MustParse("3992c8ab-6c64-4eed-b780-345d12ecc9c8"),
		CreatedBy:       createdBy,
		Message:         "It´s my first post!",
		MediaContentUrl: "",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
)

func TestTweetsDB_CreateTweetSuccess(t *testing.T) {
	conn, mock := db.NewMockPostgresConnection(t)

	repo := NewTweetRepository(conn)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "tweets" ("id","created_by","message","media_content_url","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`,
	)).WithArgs(
		tweet.ID,
		tweet.CreatedBy,
		tweet.Message,
		tweet.MediaContentUrl,
		tweet.CreatedAt,
		tweet.UpdatedAt,
		tweet.DeletedAt,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()

	_, err := repo.CreateTweet(context.Background(), tweet)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestTweetsDB_CreateTweetError(t *testing.T) {
	conn, mock := db.NewMockPostgresConnection(t)
	repo := NewTweetRepository(conn)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "tweets" ("id","created_by","message","media_content_url","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`,
	)).WithArgs(
		tweet.ID,
		tweet.CreatedBy,
		tweet.Message,
		tweet.MediaContentUrl,
		tweet.CreatedAt,
		tweet.UpdatedAt,
		tweet.DeletedAt,
	).WillReturnError(
		errors.New("error saving tweet"),
	)
	mock.ExpectRollback()

	_, err := repo.CreateTweet(context.TODO(), tweet)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
