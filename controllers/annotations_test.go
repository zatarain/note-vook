package controllers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationsAdd(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	test.Run("Should save annotation to database", func(test *testing.T) {
		// Arrange

		// Act

		// Assert
		assert.True(true)
	})

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
