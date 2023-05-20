package models

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(test *testing.T) {
	// Arrange
	assert := assert.New(test)
	now, _ := time.Parse(time.RFC1123, "Fri, 10 Jan 1986 10:04:00 GMT")
	user := &User{
		ID:        1,
		Nickname:  "dummy-user",
		Password:  "top-secret",
		CreatedAt: now,
		UpdatedAt: now.Add(7 * time.Hour),
	}
	expected := strings.Join([]string{
		"ID = 1",
		"Nickname = 'dummy-user'",
		"Created At = 'Fri, 10 Jan 1986 10:04:00 GMT'",
		"Updated At = 'Fri, 10 Jan 1986 17:04:00 GMT'",
	}, ", ")

	// Act
	actual := user.String()

	// Assert
	assert.Equal(expected, actual)
}
