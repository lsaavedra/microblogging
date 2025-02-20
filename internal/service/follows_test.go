package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microblogging/internal/eventstub"
	"microblogging/internal/model"
	"microblogging/internal/service/mocks"
)

var (
	followerID = uuid.MustParse("55f9c1bb-ff47-44f7-b05d-87967b9eed2a")
	followedID = uuid.MustParse("f2f68d10-66ad-4ae0-817d-fdea6f0d28ad")
)

func TestFollowsService_FollowUser(t *testing.T) {
	type fields struct {
		Repo      FollowersRepo
		publisher eventstub.EventProcessor
	}
	type args struct {
		ctx          context.Context
		userID       uuid.UUID
		targetUserID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test Create Follower relation ok",
			fields: fields{
				Repo:      getRepositoryMockOk(t),
				publisher: mocks.NewPubSubMock(t),
			},
			args: args{
				ctx:          context.TODO(),
				userID:       followerID,
				targetUserID: followedID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "Test Create Follower relation fails by saving into db",
			fields: fields{
				Repo:      getRepositoryMockWithFail(t),
				publisher: mocks.NewPubSubMock(t),
			},
			args: args{
				ctx:          context.TODO(),
				userID:       followerID,
				targetUserID: followedID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FollowsService{
				repo:      tt.fields.Repo,
				publisher: tt.fields.publisher,
			}
			tt.wantErr(t, s.FollowUser(tt.args.ctx, tt.args.userID, tt.args.targetUserID), fmt.Sprintf("FollowUser(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.targetUserID))
		})
	}
}

func TestFollowsService_GetFollowers(t *testing.T) {
	type fields struct {
		repo      func() FollowersRepo
		publisher eventstub.EventProcessor
	}
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Follow
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test GetFollowersOk by empty list",
			fields: fields{
				repo: func() FollowersRepo {
					repoMock := new(mocks.FollowsRepositoryMock)
					repoMock.On("GetFollowers", mock.Anything, followedID).Return([]model.Follow{}, nil)

					return repoMock
				},
				publisher: mocks.NewPubSubMock(t),
			},
			args: args{ctx: context.TODO(), userID: followedID},
			want: []model.Follow{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "Test GetFollowers fails by db error",
			fields: fields{
				repo: func() FollowersRepo {
					repoMock := new(mocks.FollowsRepositoryMock)
					repoMock.On("GetFollowers", mock.Anything, followedID).
						Return([]model.Follow{}, errors.New("error getting from db"))

					return repoMock
				},
				publisher: mocks.NewPubSubMock(t),
			},
			args: args{ctx: context.TODO(), userID: followedID},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotEmpty(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FollowsService{
				repo:      tt.fields.repo(),
				publisher: tt.fields.publisher,
			}
			got, err := s.GetFollowers(tt.args.ctx, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFollowers(%v, %v)", tt.args.ctx, tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFollowers(%v, %v)", tt.args.ctx, tt.args.userID)
		})
	}
}

func getRepositoryMockOk(t *testing.T) *mocks.FollowsRepositoryMock {
	repoMock := mocks.NewFollowsRepositoryMock(t)

	repoMock.On("CreateFollowing", mock.Anything, mock.AnythingOfType("model.Follow")).
		Return(nil)

	return repoMock
}

func getRepositoryMockWithFail(t *testing.T) *mocks.FollowsRepositoryMock {
	repoMock := mocks.NewFollowsRepositoryMock(t)

	repoMock.On("CreateFollowing", mock.Anything, mock.AnythingOfType("model.Follow")).
		Return(errors.New("error saving into DB"))

	return repoMock
}
