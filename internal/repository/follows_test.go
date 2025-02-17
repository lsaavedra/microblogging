package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"microblogging/db"
	"microblogging/internal/model"
)

func TestFollowsDB_UserFollowingAnotherUserSuccess(t *testing.T) {
	var (
		followerID  = uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
		followingID = uuid.MustParse("3992c8ab-6c64-4eed-b780-345d12ecc9c8")
	)

	conn, mock := db.NewMockPostgresConnection(t)

	repo := NewFollowsRepository(conn)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "follows" ("follower_id","following_id","followed_at") VALUES ($1,$2,$3)`,
	)).WithArgs(
		followerID,
		followingID,
		sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) // (LastInsertID, RowsAffected)
	mock.ExpectCommit()

	err := repo.SaveFollowAndFollowing(context.TODO(), model.Follow{FollowerID: followerID, FollowingID: followingID})

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowsDB_UserFollowingAnotherUserFailsByRecordsExists(t *testing.T) {
	var (
		followerID  = uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
		followingID = uuid.MustParse("3992c8ab-6c64-4eed-b780-345d12ecc9c8")
	)

	conn, mock := db.NewMockPostgresConnection(t)

	repo := NewFollowsRepository(conn)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "follows" ("follower_id","following_id","followed_at") VALUES ($1,$2,$3)`,
	)).WithArgs(
		followerID,
		followingID,
		sqlmock.AnyArg()).
		WillReturnError(errors.New("follow and following id exists"))
	mock.ExpectRollback()

	err := repo.SaveFollowAndFollowing(context.TODO(), model.Follow{FollowerID: followerID, FollowingID: followingID})

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
