package controllers

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(test *testing.T) {
	assert := assert.New(test)
	require := require.New(test)
	gin.SetMode(gin.TestMode)
	server := gin.New()
	server.HEAD("/health", HealthCheck)
	recorder := httptest.NewRecorder()

	request, exception := http.NewRequest(http.MethodHead, "/health", nil)
	require.Nil(exception)

	// Perform the request
	server.ServeHTTP(recorder, request)

	// Check to see if the response was what you expected
	assert.Equal(http.StatusOK, recorder.Code)
}
