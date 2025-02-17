package handler

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"microblogging/db"
	"microblogging/internal/model"
	"microblogging/testharness"
	"microblogging/testharness/utils"
)

func TestAddNewFollowAndFollowingRelation(t *testing.T) {
	database := testharness.NewDbForTest()

	tx := database.Begin()
	defer tx.Rollback()

	router := routerSetup(tx)

	tests := []struct {
		name        string
		body        model.UserFollowingRequestBody
		expectedLen int
	}{
		{
			"test user is following a new one ok",
			model.UserFollowingRequestBody{
				UserID:       uuid.MustParse("b9a9f007-0b17-46ed-a165-a49d28238bbb"),
				TargetUserID: uuid.MustParse("3992c8ab-6c64-4eed-b780-345d12ecc9c8"),
			},
			1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := utils.NewRequestWithResponse[any, any](router)
			response := req.
				URL("/api/v1/users/follow").
				POST(test.body).
				Expect(http.StatusOK)

			assert.Nil(t, response)

			reqTwo := utils.NewRequestWithResponse[any, []model.Follow](router)
			result := reqTwo.
				URL("/api/v1/users/follow").
				GET().
				Expect(http.StatusOK).
				Body()

			assert.Equal(t, response, len(result))
		})
	}
}

func routerSetup(tx *gorm.DB) *gin.Engine {
	return RouterWithHandlers(
		&db.Database{DB: tx},
	)
}
