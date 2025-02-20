package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microblogging/internal/eventstub"
	"microblogging/internal/model"
	"microblogging/internal/service/mocks"
)

func TestTweetsService_CreateUserTweet(t1 *testing.T) {
	var tweetCreatedID = uuid.New()
	type fields struct {
		logger    *log.Logger
		repo      func() TweetsRepo
		publisher eventstub.EventProcessor
	}
	type args struct {
		ctx          context.Context
		tweetRequest model.CreateTweetRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.CreateTweetResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test Create user tweet ok",
			fields: fields{
				logger: log.New(),
				repo: func() TweetsRepo {
					repoMock := new(mocks.TweetsRepositoryMock)
					repoMock.On("CreateTweet", mock.Anything, mock.AnythingOfType("model.Tweet")).
						Return(model.Tweet{ID: tweetCreatedID}, nil)

					return repoMock
				},
				publisher: new(mocks.PubSubMock),
			},
			args: args{
				ctx: context.TODO(),
				tweetRequest: model.CreateTweetRequest{
					UserID:  followedID,
					Message: "This is my first tweet",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Empty(t, err)
			},
		},
		{
			name: "Test Create user tweet fails by getting values from db",
			fields: fields{
				logger: log.New(),
				repo: func() TweetsRepo {
					repoMock := new(mocks.TweetsRepositoryMock)
					repoMock.On("CreateTweet", mock.Anything, mock.AnythingOfType("model.Tweet")).
						Return(model.Tweet{}, errors.New("error getting from db"))

					return repoMock
				},
				publisher: new(mocks.PubSubMock),
			},
			args: args{
				ctx: context.TODO(),
				tweetRequest: model.CreateTweetRequest{
					UserID:  followedID,
					Message: "This is my first tweet",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotEmpty(t, err)
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TweetsService{
				logger:    tt.fields.logger,
				repo:      tt.fields.repo(),
				publisher: tt.fields.publisher,
			}
			_, err := t.CreateUserTweet(tt.args.ctx, tt.args.tweetRequest)
			if !tt.wantErr(t1, err, fmt.Sprintf("CreateUserTweet(%v, %v)", tt.args.ctx, tt.args.tweetRequest)) {
				return
			}
		})
	}
}
