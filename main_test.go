package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(test *testing.T) {
	assert := assert.New(test)

	// Teardown test suite
	defer log.SetOutput(os.Stderr)

	test.Run("Should print project name", func(test *testing.T) {
		// Arrange
		var capture bytes.Buffer
		log.SetOutput(&capture)

		// Act
		main()

		// Assert
		assert.Contains(capture.String(), "NoteVook")
	})
}
