package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
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
			database.
				On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", video.ID, current.ID).
				Return(&gorm.DB{Error: nil}).Run(
				func(arguments mock.Arguments) {
					result := arguments.Get(0).(*models.Video)
					*result = video
				},
			)

			database.
				On("Create", mock.AnythingOfType("*models.Annotation")).
				Return(&gorm.DB{Error: nil}).Run(
				func(arguments mock.Arguments) {
					recordset := arguments.Get(0).(*models.Annotation)
					*recordset = testcase.Expected
				},
			)

			server.POST("/annotations", authorise(&current), annotations.Add)
			body, _ := json.Marshal(&testcase.Input)
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

	invalidInputs := []struct {
		Input    gin.H
		Expected string
	}{
		{
			Input: gin.H{
				"video_id": video.ID,
				"start":    7*60 + 20,
				"end":      "07:30",
			},
			Expected: "Field validation for 'Title' failed on the 'required' tag",
		},
		{
			Input: gin.H{
				"video_id": "fake.ID",
				"type":     1,
				"title":    "My dummy annotation",
				"notes":    "My additional notes",
				"start":    "07:28",
				"end":      "07:30",
			},
			Expected: "cannot unmarshal string into Go struct field AddAnnotationContract.video_id of type uint",
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     "wrong",
				"title":    "My dummy annotation",
				"notes":    "My additional notes",
				"start":    10,
				"end":      30,
			},
			Expected: "cannot unmarshal string into Go struct field AddAnnotationContract.type of type uint",
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     1,
				"title":    "My dummy annotation",
				"start":    "7-45",
				"end":      "07m30s",
			},
			Expected: "time: unknown unit",
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     1,
				"title":    "My dummy annotation",
				"start":    "7:28",
				"end":      "7-30",
			},
			Expected: "time: unknown unit",
		},
		{
			Input: gin.H{
				"video_id": video.ID,
				"type":     1,
				"title":    "My dummy annotation",
				"start":    "01:25",
				"end":      "01:10",
			},
			Expected: "Field validation for 'Start' failed on the 'ltefield' tag",
		},
	}

	for _, testcase := range invalidInputs {
		test.Run("Should NOT save annotation on invalid body input", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			annotations := &AnnotationsController{Database: database}

			server.POST("/annotations", authorise(&current), annotations.Add)
			body, _ := json.Marshal(&testcase.Input)
			request, _ := http.NewRequest(http.MethodPost, "/annotations", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Contains(recorder.Body.String(), "Failed to read input")
			assert.Contains(recorder.Body.String(), testcase.Expected)
			database.AssertExpectations(test)
		})
	}

	test.Run("Should NOT save annotation if video doesn't exists for current user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		annotations := &AnnotationsController{Database: database}
		database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", video.ID, current.ID).
			Return(&gorm.DB{Error: errors.New("invalid video_id")}).Run(
			func(arguments mock.Arguments) {
				result := arguments.Get(0).(*models.Video)
				*result = video
			},
		)

		server.POST("/annotations", authorise(&current), annotations.Add)
		body, _ := json.Marshal(&gin.H{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"start":    "07m28s",
			"end":      "07m30s",
		})
		request, _ := http.NewRequest(http.MethodPost, "/annotations", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusNotFound, recorder.Code)
		assert.Contains(recorder.Body.String(), "Video not found")
		assert.Contains(recorder.Body.String(), "invalid video_id")
		database.AssertExpectations(test)
	})

	invalidTimeRanges := []gin.H{
		{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"start":    video.Duration + 7,
			"end":      video.Duration + 10,
		},
		{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"start":    7,
			"end":      video.Duration + 10,
		},
		{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"start":    -7,
			"end":      video.Duration + 10,
		},
		{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"start":    -17,
			"end":      video.Duration - 10,
		},
	}

	for _, testcase := range invalidTimeRanges {
		test.Run("Should NOT save annotation either Start or End are out of bounds of video Duration", func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			annotations := &AnnotationsController{Database: database}
			database.
				On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", video.ID, current.ID).
				Return(&gorm.DB{Error: nil}).Run(
				func(arguments mock.Arguments) {
					result := arguments.Get(0).(*models.Video)
					*result = video
				},
			)

			server.POST("/annotations", authorise(&current), annotations.Add)
			body, _ := json.Marshal(&testcase)
			request, _ := http.NewRequest(http.MethodPost, "/annotations", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Contains(recorder.Body.String(), "Invalid time interval")
			assert.Contains(recorder.Body.String(), "start and end must be positive and less or equal than video duration")
			database.AssertExpectations(test)
		})
	}

	test.Run("Should response with an error when database returns error", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		annotations := &AnnotationsController{Database: database}
		database.
			On("First", mock.AnythingOfType("*models.Video"), "id = ? AND user_id = ?", video.ID, current.ID).
			Return(&gorm.DB{Error: nil}).Run(
			func(arguments mock.Arguments) {
				result := arguments.Get(0).(*models.Video)
				*result = video
			},
		)

		database.
			On("Create", mock.AnythingOfType("*models.Annotation")).
			Return(&gorm.DB{Error: errors.New("database insertion error")})

		server.POST("/annotations", authorise(&current), annotations.Add)
		body, _ := json.Marshal(&gin.H{
			"video_id": video.ID,
			"type":     1,
			"title":    "My dummy annotation",
			"notes":    "My additional notes",
			"start":    "07:28",
			"end":      "07:30",
		})
		request, _ := http.NewRequest(http.MethodPost, "/annotations", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to save the annotation")
		assert.Contains(recorder.Body.String(), "database insertion error")
		database.AssertExpectations(test)
	})
}

func TestAnnotationsDelete(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	date, _ := time.Parse(time.DateOnly, "2021-01-01")
	current := models.User{
		ID:       3,
		Nickname: "dummy",
	}

	annotation := &models.Annotation{
		ID:        12,
		VideoID:   7,
		Type:      0,
		Title:     "My annotation",
		Notes:     "",
		Start:     15,
		End:       30,
		CreatedAt: date,
		UpdatedAt: date.Add(40 * time.Hour),
	}
	test.Run("Should delete the annotation", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		annotations := &AnnotationsController{Database: database}

		gormFakeSuccess := &gorm.DB{Error: nil}
		database.On("Joins", "Video").Return(gormFakeSuccess)
		var arguments struct {
			ValueType  string
			Conditions []interface{}
		}
		monkey.PatchInstanceMethod(
			reflect.TypeOf(gormFakeSuccess),
			"First",
			func(DB *gorm.DB, value interface{}, conditions ...interface{}) *gorm.DB {
				recordset := value.(*models.Annotation)
				*recordset = *annotation
				arguments.ValueType = reflect.TypeOf(value).String()
				arguments.Conditions = conditions
				return gormFakeSuccess
			},
		)
		defer monkey.UnpatchAll()

		database.On("Delete", mock.AnythingOfType("*models.Annotation")).Return(gormFakeSuccess)

		server.DELETE("/annotations/:id", authorise(&current), annotations.Delete)
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/annotations/%d", annotation.ID), nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusOK, recorder.Code)
		assert.Contains(recorder.Body.String(), "Annotation successfully deleted")

		assert.Equal("*models.Annotation", arguments.ValueType)
		assert.Len(arguments.Conditions, 3)
		assert.Equal("annotations.id = ? AND user_id = ?", arguments.Conditions[0])
		assert.Equal(fmt.Sprint(annotation.ID), arguments.Conditions[1])
		assert.Equal(current.ID, arguments.Conditions[2])
		database.AssertExpectations(test)
	})

	test.Run("Should response with HTTP 404 when it doesn't exits in the database", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		annotations := &AnnotationsController{Database: database}

		gormFakeSuccess := &gorm.DB{Error: nil}
		database.On("Joins", "Video").Return(gormFakeSuccess)
		var arguments struct {
			ValueType  string
			Conditions []interface{}
		}
		monkey.PatchInstanceMethod(
			reflect.TypeOf(gormFakeSuccess),
			"First",
			func(DB *gorm.DB, value interface{}, conditions ...interface{}) *gorm.DB {
				arguments.ValueType = reflect.TypeOf(value).String()
				arguments.Conditions = conditions
				return &gorm.DB{Error: errors.New("no results")}
			},
		)
		defer monkey.UnpatchAll()

		server.DELETE("/annotations/:id", authorise(&current), annotations.Delete)
		request, _ := http.NewRequest(http.MethodDelete, "/annotations/95", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusNotFound, recorder.Code)
		assert.Contains(recorder.Body.String(), "Annotation not found")

		assert.Equal("*models.Annotation", arguments.ValueType)
		assert.Len(arguments.Conditions, 3)
		assert.Equal("annotations.id = ? AND user_id = ?", arguments.Conditions[0])
		assert.Equal("95", arguments.Conditions[1])
		assert.Equal(current.ID, arguments.Conditions[2])
		database.AssertExpectations(test)
	})

	test.Run("Should response with HTTP 400 when the database returns an error", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		annotations := &AnnotationsController{Database: database}

		gormFakeSuccess := &gorm.DB{Error: nil}
		database.On("Joins", "Video").Return(gormFakeSuccess)
		var arguments struct {
			ValueType  string
			Conditions []interface{}
		}
		monkey.PatchInstanceMethod(
			reflect.TypeOf(gormFakeSuccess),
			"First",
			func(DB *gorm.DB, value interface{}, conditions ...interface{}) *gorm.DB {
				recordset := value.(*models.Annotation)
				*recordset = *annotation
				arguments.ValueType = reflect.TypeOf(value).String()
				arguments.Conditions = conditions
				return gormFakeSuccess
			},
		)
		defer monkey.UnpatchAll()

		database.On("Delete", mock.AnythingOfType("*models.Annotation")).
			Return(&gorm.DB{Error: errors.New("unable to delete")})

		server.DELETE("/annotations/:id", authorise(&current), annotations.Delete)
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/annotations/%d", annotation.ID), nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to delete the annotation")
		assert.Contains(recorder.Body.String(), "unable to delete")

		assert.Equal("*models.Annotation", arguments.ValueType)
		assert.Len(arguments.Conditions, 3)
		assert.Equal("annotations.id = ? AND user_id = ?", arguments.Conditions[0])
		assert.Equal(fmt.Sprint(annotation.ID), arguments.Conditions[1])
		assert.Equal(current.ID, arguments.Conditions[2])
		database.AssertExpectations(test)
	})
}
