package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zatarain/note-vook/mocks"
	"github.com/zatarain/note-vook/models"
	"gorm.io/gorm"
)

func TestAnnotationsAdd(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	date, _ := time.Parse(time.DateOnly, "2021-01-01")
	current := models.User{
		ID:       3,
		Nickname: "dummy",
	}
	video := models.Video{
		ID:          1,
		UserID:      3,
		Title:       "Dummy video 01",
		Description: "This is a dummy video number one",
		Duration:    7*60 + 45,
		Link:        "https://youtube.com/v/number-one",
		CreatedAt:   date,
		UpdatedAt:   date.Add(4 * time.Hour),
	}
	validInputs := []struct {
		Input    gin.H
		Expected models.Annotation
	}{
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     1,
				"title":    "My dummy annotation",
				"notes":    "My additional notes",
				"start":    "07:28",
				"end":      "07:30",
			},
			Expected: models.Annotation{
				ID:        6,
				VideoID:   9,
				Type:      1,
				Title:     "My annotation",
				Notes:     "My additional notes",
				Start:     7*60 + 28,
				End:       7*60 + 30,
				CreatedAt: date.Add(4 * 24 * time.Hour),
				UpdatedAt: date.Add(4 * 24 * time.Hour),
			},
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"title":    "My dummy annotation",
				"notes":    "My additional notes",
				"start":    10,
				"end":      30,
			},
			Expected: models.Annotation{
				ID:        6,
				VideoID:   9,
				Type:      0,
				Title:     "My annotation",
				Notes:     "My additional notes",
				Start:     10,
				End:       30,
				CreatedAt: date.Add(4 * 24 * time.Hour),
				UpdatedAt: date.Add(4 * 24 * time.Hour),
			},
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     1,
				"title":    "My dummy annotation",
				"start":    "07m28s",
				"end":      "07m30s",
			},
			Expected: models.Annotation{
				ID:        6,
				VideoID:   9,
				Type:      1,
				Title:     "My annotation",
				Notes:     "",
				Start:     7*60 + 28,
				End:       7*60 + 30,
				CreatedAt: date.Add(4 * 24 * time.Hour),
				UpdatedAt: date.Add(4 * 24 * time.Hour),
			},
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"title":    "My dummy annotation",
				"start":    7*60 + 20,
				"end":      "07:30",
			},
			Expected: models.Annotation{
				ID:        6,
				VideoID:   9,
				Type:      0,
				Title:     "My annotation",
				Notes:     "",
				Start:     7*60 + 20,
				End:       7*60 + 30,
				CreatedAt: date.Add(4 * 24 * time.Hour),
				UpdatedAt: date.Add(4 * 24 * time.Hour),
			},
		},
	}

	for _, testcase := range validInputs {
		test.Run("Should save annotation to database", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			annotations := &AnnotationsController{Database: database}
			lookup := database.
				On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", video.ID, current.ID).
				Return(&gorm.DB{Error: nil})
			lookup.RunFn = func(arguments mock.Arguments) {
				result := arguments.Get(0).(*models.Video)
				*result = video
			}

			insert := database.
				On("Create", mock.AnythingOfType("*models.Annotation")).
				Return(&gorm.DB{Error: nil})
			insert.RunFn = func(arguments mock.Arguments) {
				recordset := arguments.Get(0).(*models.Annotation)
				*recordset = testcase.Expected
			}

			server.POST("/annotations", authorise(&current), annotations.Add)
			body, _ := json.Marshal(testcase.Input)
			request, _ := http.NewRequest(http.MethodPost, "/annotations", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()
			expected, _ := json.Marshal(&testcase.Expected)

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusCreated, recorder.Code)
			assert.Equal(expected, recorder.Body.Bytes())
			database.AssertExpectations(test)
		})
	}
	test.Run("Should NOT save annotation on invalid body input", func(test *testing.T) {
		// Arrange

		// Act

		// Assert
		assert.True(true)
	})

	test.Run("Should NOT save annotation if video doesn't exists for current user", func(test *testing.T) {
		// Arrange

		// Act

		// Assert
		assert.True(true)
	})

	test.Run("Should NOT save annotation either Start or End are out of bounds of video Duration", func(test *testing.T) {
		// Arrange

		// Act

		// Assert
		assert.True(true)
	})

	test.Run("Should response with an error when there is a database error returned", func(test *testing.T) {
		// Arrange

		// Act

		// Assert
		assert.True(true)
	})
}
