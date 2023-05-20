package controllers

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetVideos(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	server := gin.New()
	server.HEAD("/videos", GetVideos)
	recorder := httptest.NewRecorder()

	request, exception := http.NewRequest(http.MethodHead, "/videos", nil)
	assert.Nil(exception)

	// Perform the request
	server.ServeHTTP(recorder, request)

	// Check to see if the response was what you expected
	assert.Equal(http.StatusOK, recorder.Code)
}
