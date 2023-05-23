package configuration

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/zatarain/note-vook/mocks"
)

func TestSetup(test *testing.T) {
	test.Run("Should setup all the end-points", func(test *testing.T) {
		// Arrange
		server := new(mocks.MockedEngine)
		endPointHandler := mock.AnythingOfType("gin.HandlerFunc")
		authorisationHandler := mock.AnythingOfType("gin.HandlerFunc")
		server.On("HEAD", "/health", endPointHandler).Return(server)
		server.On("POST", "/signup", endPointHandler).Return(server)
		server.On("POST", "/login", endPointHandler).Return(server)

		// Authorised end-points
		server.On("GET", "/videos", authorisationHandler, endPointHandler).Return(server)
		server.On("POST", "/videos", authorisationHandler, endPointHandler).Return(server)
		server.On("GET", "/videos/:id", authorisationHandler, endPointHandler).Return(server)
		server.On("PATCH", "/videos/:id", authorisationHandler, endPointHandler).Return(server)
		server.On("DELETE", "/videos/:id", authorisationHandler, endPointHandler).Return(server)

		server.On("POST", "/annotations", authorisationHandler, endPointHandler).Return(server)
		server.On("PATCH", "/annotations/:id", authorisationHandler, endPointHandler).Return(server)
		server.On("DELETE", "/annotations/:id", authorisationHandler, endPointHandler).Return(server)

		// Act
		Setup(server)

		// Assert
		server.AssertExpectations(test)
	})
}
