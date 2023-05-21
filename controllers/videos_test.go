package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func authorise(current *models.User) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("user", current)
	}
}

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
			CreatedAt:   dummyDate.Add(-7 * time.Hour),
			UpdatedAt:   dummyDate,
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

func TestVideosView(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	dummyDate, _ := time.Parse(time.DateOnly, "2021-01-01")
	video := models.Video{
		ID:          3,
		UserID:      3,
		Title:       "Dummy video 03",
		Description: "This is a dummy video number three",
		Duration:    50,
		Link:        "https://youtube.com/v/number-three",
		CreatedAt:   dummyDate,
		UpdatedAt:   dummyDate,
	}
	current := models.User{
		ID:       3,
		Nickname: "three",
	}

	test.Run("Should return the video for the current user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		call := database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", fmt.Sprint(video.ID), current.ID).
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			recordset := arguments.Get(0).(*models.Video)
			*recordset = video
		}
		server.GET("/videos/:id", authorise(&current), videos.View)
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/videos/%d", video.ID), nil)
		recorder := httptest.NewRecorder()
		expected, _ := json.Marshal(&video)

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusOK, recorder.Code)
		assert.Equal(expected, recorder.Body.Bytes())
		database.AssertExpectations(test)
	})

	test.Run("Should return HTTP 404 if video it's not in database", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", "10", current.ID).
			Return(&gorm.DB{Error: errors.New("no results")})
		server.GET("/videos/:id", authorise(&current), videos.View)
		request, _ := http.NewRequest(http.MethodGet, "/videos/10", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusNotFound, recorder.Code)
		assert.Contains(recorder.Body.String(), "Video not found")
		database.AssertExpectations(test)
	})
}

func TestVideosAdd(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	dummyDate, _ := time.Parse(time.DateOnly, "2021-01-01")
	current := models.User{
		ID:       3,
		Nickname: "three",
	}

	validTestcases := []struct {
		Input    gin.H
		Expected models.Video
	}{
		{
			Expected: models.Video{
				ID:          3,
				UserID:      3,
				Title:       "Dummy video 03",
				Description: "This is a dummy video number three",
				Duration:    105,
				Link:        "https://youtube.com/v/number-three",
				CreatedAt:   dummyDate,
				UpdatedAt:   dummyDate,
			},
			Input: gin.H{
				"title":       "Dummy video 03",
				"description": "This is a dummy video number three",
				"duration":    "1:45",
				"link":        "https://youtube.com/v/number-three",
			},
		},
		{
			Expected: models.Video{
				ID:          5,
				UserID:      3,
				Title:       "Dummy video 05",
				Description: "",
				Duration:    400,
				Link:        "https://youtube.com/v/number-five",
				CreatedAt:   dummyDate,
				UpdatedAt:   dummyDate,
			},
			Input: gin.H{
				"title":    "Dummy video 04",
				"duration": 400,
				"link":     "https://youtube.com/v/number-five",
			},
		},
	}

	for _, testcase := range validTestcases {
		test.Run("Should create a new video owned by current user when valid data is received", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			videos := &VideosController{Database: database}
			call := database.
				On("Create", mock.AnythingOfType("*models.Video")).
				Return(&gorm.DB{Error: nil})
			call.RunFn = func(arguments mock.Arguments) {
				recordset := arguments.Get(0).(*models.Video)
				*recordset = testcase.Expected
			}
			server.POST("/videos", authorise(&current), videos.Add)

			body, _ := json.Marshal(testcase.Input)
			request, _ := http.NewRequest(http.MethodPost, "/videos", bytes.NewBuffer(body))
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

	invalidTestcases := []struct {
		Expected string
		Input    gin.H
	}{
		{
			Expected: "Field validation for 'Title' failed on the 'required' tag",
			Input: gin.H{
				"description": "This is a dummy video number three",
				"duration":    "1:45",
				"link":        "https://youtube.com/v/number-three",
			},
		},
		{
			Expected: "Field validation for 'Title' failed on the 'required' tag",
			Input: gin.H{
				"title":       "",
				"description": "This is a dummy video number three",
				"duration":    "1:45",
				"link":        "https://youtube.com/v/number-three",
			},
		},
		{
			Expected: "Field validation for 'Link' failed on the 'required' tag",
			Input: gin.H{
				"title":       "Dummy video 03",
				"description": "This is a dummy video number three",
				"duration":    "1:45",
			},
		},
		{
			Expected: "Field validation for 'Duration' failed on the 'required' tag",
			Input: gin.H{
				"title":       "Dummy video 03",
				"description": "This is a dummy video number three",
				"link":        "https://youtube.com/v/number-three",
			},
		},
		{
			Expected: "Field validation for 'Link' failed on the 'url' tag",
			Input: gin.H{
				"title":       "Dummy video 03",
				"description": "This is a dummy video number three",
				"duration":    "1:45",
				"link":        "bad link address",
			},
		},
		{
			Expected: "time: invalid duration",
			Input: gin.H{
				"title":       "Dummy video 03",
				"description": "This is a dummy video number three",
				"duration":    "hello",
				"link":        "https://youtube.com/v/number-three",
			},
		},
	}

	for _, testcase := range invalidTestcases {
		test.Run("Should NOT try to create for current user with invalid inputs", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			videos := &VideosController{Database: database}
			database.
				On("Create", mock.AnythingOfType("*models.Video")).
				Return(&gorm.DB{Error: nil})
			server.POST("/videos", authorise(&current), videos.Add)

			body, _ := json.Marshal(testcase.Input)
			request, _ := http.NewRequest(http.MethodPost, "/videos", bytes.NewBuffer([]byte(body)))
			recorder := httptest.NewRecorder()

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Contains(recorder.Body.String(), testcase.Expected)
			assert.Contains(recorder.Body.String(), "Failed to read input")
			database.AssertNotCalled(test, "Create", mock.AnythingOfType("*models.Video"))
		})
	}

	test.Run("Should response HTTP 400 when there is a problem with database (e. g. unique index)", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.Video")).
			Return(&gorm.DB{Error: errors.New("unique index violation")})
		server.POST("/videos", authorise(&current), videos.Add)

		body, _ := json.Marshal(validTestcases[0].Input)
		request, _ := http.NewRequest(http.MethodPost, "/videos", bytes.NewBuffer([]byte(body)))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to save the video")
		assert.Contains(recorder.Body.String(), "unique index violation")
		database.AssertExpectations(test)
	})
}

func TestVideosEdit(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	/**
	dummyDate, _ := time.Parse(time.DateOnly, "2021-01-01")
	video := models.Video{
		ID:          3,
		UserID:      3,
		Title:       "Dummy video 03",
		Description: "This is a dummy video number three",
		Duration:    50,
		Link:        "https://youtube.com/v/number-three",
		CreatedAt:   dummyDate,
		UpdatedAt:   dummyDate,
	}
	/**/
	current := models.User{
		ID:       3,
		Nickname: "three",
	}

	test.Run("Should return HTTP 404 if video it's not in database", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", "75", current.ID).
			Return(&gorm.DB{Error: errors.New("no results")})
		server.PATCH("/videos/:id", authorise(&current), videos.Edit)
		request, _ := http.NewRequest(http.MethodPatch, "/videos/75", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusNotFound, recorder.Code)
		assert.Contains(recorder.Body.String(), "Video not found")
		database.AssertExpectations(test)
	})
}

func TestVideosDelete(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	/**
	dummyDate, _ := time.Parse(time.DateOnly, "2021-01-01")
	video := models.Video{
		ID:          3,
		UserID:      3,
		Title:       "Dummy video 03",
		Description: "This is a dummy video number three",
		Duration:    50,
		Link:        "https://youtube.com/v/number-three",
		CreatedAt:   dummyDate,
		UpdatedAt:   dummyDate,
	}
	/**/
	current := models.User{
		ID:       3,
		Nickname: "three",
	}

	test.Run("Should return HTTP 404 if video it's not in database", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		videos := &VideosController{Database: database}
		database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", "32", current.ID).
			Return(&gorm.DB{Error: errors.New("no results")})
		server.DELETE("/videos/:id", authorise(&current), videos.Edit)
		request, _ := http.NewRequest(http.MethodDelete, "/videos/32", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusNotFound, recorder.Code)
		assert.Contains(recorder.Body.String(), "Video not found")
		database.AssertExpectations(test)
	})
}
