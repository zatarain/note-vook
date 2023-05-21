package controllers

import (
	"encoding/json"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zatarain/note-vook/mocks"
	"github.com/zatarain/note-vook/models"
	"gorm.io/gorm"
)

func TestVideosIndex(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	dummyDate, _ := time.Parse(time.DateOnly, "2021-01-01")
	dataset := []models.Video{
		{
			ID:          1,
			UserID:      3,
			Title:       "Dummy video 01",
			Description: "This is a dummy video number one",
			Duration:    100,
			Link:        "https://youtube.com/v/number-one",
			CreatedAt:   dummyDate,
			UpdatedAt:   dummyDate.Add(4 * time.Hour),
		},
		{
			ID:          2,
			UserID:      4,
			Title:       "Dummy video 02",
			Description: "This is a dummy video number two",
			Duration:    200,
			Link:        "https://youtube.com/v/number-two",
			CreatedAt:   dummyDate,
			UpdatedAt:   dummyDate.Add(7 * time.Hour),
		},
		{
			ID:          3,
			UserID:      3,
			Title:       "Dummy video 03",
			Description: "This is a dummy video number three",
			Duration:    50,
			Link:        "https://youtube.com/v/number-three",
			CreatedAt:   dummyDate,
			UpdatedAt:   dummyDate,
		},
	}

	users := []models.User{
		{ID: 3, Nickname: "three"},
		{ID: 4, Nickname: "four"},
		{ID: 5, Nickname: "five"},
	}

	authorise := func(current *models.User) gin.HandlerFunc {
		return func(context *gin.Context) {
			context.Set("user", current)
		}
	}

	resultsets := map[int][]models.Video{
		3: {dataset[0], dataset[2]},
		4: {dataset[1]},
		5: {},
	}

	for _, current := range users {
		test.Run("Should return the list of videos for the current user", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			videos := &VideosController{Database: database}
			call := database.
				On("Find", mock.AnythingOfType("*[]models.Video"), "user_id = ?", current.ID).
				Return(&gorm.DB{Error: nil})
			call.RunFn = func(arguments mock.Arguments) {
				recordset := arguments.Get(0).(*[]models.Video)
				*recordset = resultsets[current.ID]
				//call.ReturnArguments = mock.Arguments{db, *recordset}
			}
			server.GET("/videos", authorise(&current), videos.Index)
			request, _ := http.NewRequest(http.MethodGet, "/videos", nil)
			recorder := httptest.NewRecorder()
			expected, _ := json.Marshal(resultsets[current.ID])

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusOK, recorder.Code)
			assert.Equal(expected, recorder.Body.Bytes())
			database.AssertExpectations(test)
		})
	}
}

/**
func TestVideosAdd(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should create a new video owned by current user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.Video")).
			Return(&gorm.DB{Error: nil})
		server.POST("/videos", videos.Add)
		video := models.Video{
			Link:     "dummy-user",
			Duration: 44,
		}
		body, _ := json.Marshal(video)
		request, _ := http.NewRequest(http.MethodPost, "/videos", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusCreated, recorder.Code)
		assert.Contains(recorder.Body.String(), "Video successfully created")
		database.AssertExpectations(test)
	})
}
/**/
