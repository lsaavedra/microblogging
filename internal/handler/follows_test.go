package handler

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"microblogging/db"
	"microblogging/internal/eventstub"
	"microblogging/internal/model"
	"microblogging/testhelpers"
	"microblogging/testhelpers/utils"
)

func TestAddNewFollowAndFollowingRelation(t *testing.T) {
	database := testhelpers.NewDbForTest()

	tx := database.Begin()
	defer tx.Rollback()

	router := RouterWithHandlers(
		&db.Database{DB: tx},
		log.New(),
		eventstub.NewPubSub(),
	)

	tests := []struct {
		name string
		body model.UserFollowingRequestBody
	}{
		{
			"test user is following a new one fails by same user",
			model.UserFollowingRequestBody{
				UserID:       uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb"),
				TargetUserID: uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := utils.NewRequestWithResponse[any, any](router)
			_ = req.
				URL("/api/v1/users/follow?user_id=b9a9f007-0b17-46ed-a165-a49d28238bbb").
				POST(test.body).
				Expect(http.StatusBadRequest)
		})
	}
}
