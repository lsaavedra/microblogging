package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"microblogging/db"
	"microblogging/internal/eventstub"
	"microblogging/internal/eventstub/mock"
	"microblogging/internal/model"
	"microblogging/testhelpers"
	"microblogging/testhelpers/utils"
)

func TestUserCreateTweetOk(t *testing.T) {
	database := testhelpers.NewDbForTest()

	tx := database.Begin()
	defer tx.Rollback()

	router := RouterWithHandlers(
		&db.Database{DB: tx},
		log.New(),
		&mock.PubSubMock{},
	)

	userID := uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
	body := model.CreateTweetRequest{
		UserID:  userID,
		Message: "First Test Tweet",
	}

	req := utils.NewRequestWithResponse[any, any](router)
	response := req.
		URL(fmt.Sprintf("/api/v1/tweets?user_id=%s", userID)).
		POST(body).
		Expect(http.StatusOK)

	assert.NotNil(t, response)

	var createdTweet model.Tweet
	err := tx.Where("created_by = ?", userID).Find(&createdTweet).Error
	if err != nil {
		panic(err)
	}

	assert.Equal(t, createdTweet.Message, body.Message)
}

func TestUserCreateTweetFailsByMaxSize(t *testing.T) {
	database := testhelpers.NewDbForTest()

	tx := database.Begin()
	defer tx.Rollback()

	router := RouterWithHandlers(
		&db.Database{DB: tx},
		log.New(),
		&mock.PubSubMock{},
	)

	userID := uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
	body := model.CreateTweetRequest{
		UserID: userID,
		Message: "test message for twitter .........................................................................." +
			"........................................................................................................" +
			"..........................................................................................................",
	}

	req := utils.NewRequestWithResponse[any, any](router)
	response := req.
		URL(fmt.Sprintf("/api/v1/tweets?user_id=%s", userID)).
		POST(body).
		Expect(http.StatusBadRequest).BodyBytes()

	assert.Contains(t, string(response), "message is too long")
}

func TestGetUserTimeline(t *testing.T) {
	database := testhelpers.NewDbForTest()

	tx := database.Begin()
	defer tx.Rollback()

	router := RouterWithHandlers(
		&db.Database{DB: tx},
		log.New(),
		eventstub.NewPubSub(),
	)

	userID := uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb")
	followedUserID := uuid.MustParse("33b9acdf-383d-4014-baf6-f9fa844a587c")

	//Following a new user
	followUserBody := model.UserFollowingRequestBody{
		UserID:       userID,
		TargetUserID: followedUserID,
	}
	req := utils.NewRequestWithResponse[any, any](router)
	response := req.
		URL(fmt.Sprintf("/api/v1/users/follow?user_id=%s", userID)).
		POST(followUserBody).
		Expect(http.StatusOK)

	assert.NotNil(t, response)

	//Get Timeline for user (empty)
	result := req.
		URL(fmt.Sprintf("/api/v1/tweets/timelines?user_id=%s", userID)).
		GET().
		Expect(http.StatusOK).
		BodyBytes()

	assert.Equal(t, "{\"total\":0,\"items\":[]}", string(result))

	// Followed user creates a new tweet
	body := model.CreateTweetRequest{
		UserID:  followedUserID,
		Message: "This is my first tweet",
	}

	response = req.
		URL(fmt.Sprintf("/api/v1/tweets?user_id=%s", followedUserID)).
		POST(body).
		Expect(http.StatusOK)

	assert.NotNil(t, response)

	time.Sleep(2 * time.Second)

	// Get timeline for user (should has data)
	result = req.
		URL(fmt.Sprintf("/api/v1/tweets/timelines?user_id=%s", userID)).
		GET().
		Expect(http.StatusOK).
		BodyBytes()

	assert.NotEmpty(t, result)
}
