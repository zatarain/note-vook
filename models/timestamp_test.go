package models

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(test *testing.T) {
	assert := assert.New(test)

	// Teardown test suite
	defer monkey.UnpatchAll()

	numbers := []struct {
		Input    string
		Expected TimeStamp
	}{
		{Input: "1500", Expected: 1500},
		{Input: "3.54", Expected: 3},
		{Input: "600", Expected: 600},
	}
	for _, testcase := range numbers {
		test.Run(fmt.Sprintf("Should convert to %v any numeric value without error from %s", testcase.Expected, testcase.Input), func(test *testing.T) {
			// Arrange
			bytes := []byte(testcase.Input)
			actual := TimeStamp(0)
			timestamp := &actual

			// Act
			exception := timestamp.UnmarshalJSON(bytes)

			// Assert
			assert.Nil(exception)
			assert.Equal(testcase.Expected, actual)
		})
	}

	strings := []struct {
		Input    string
		Expected TimeStamp
	}{
		{Input: `"1500"`, Expected: 1500},
		{Input: `"0900"`, Expected: 900},
		{Input: `"15:00"`, Expected: 15 * 60},
		{Input: `"10:45:15"`, Expected: 10*3600 + 45*60 + 15},
		{Input: `"3h4m5s"`, Expected: 3*3600 + 4*60 + 5},
		{Input: `"1:2:3"`, Expected: 3600 + 2*60 + 3},
		{Input: `"09:09:09"`, Expected: 9*3600 + 9*60 + 9},
		{Input: `"3.54"`, Expected: 3},
		{Input: `"7m15s"`, Expected: 7*60 + 15},
		{Input: `"3h4m5s"`, Expected: 3*3600 + 4*60 + 5},
		{Input: `"32h16m8s"`, Expected: 32*3600 + 16*60 + 8},
	}
	for _, testcase := range strings {
		test.Run(fmt.Sprintf("Should convert a valid quoted string value to %v without error from %s", testcase.Expected, testcase.Input), func(test *testing.T) {
			// Arrange
			bytes := []byte(testcase.Input)
			actual := TimeStamp(0)
			timestamp := &actual

			// Act
			exception := timestamp.UnmarshalJSON(bytes)

			// Assert
			assert.Nil(exception)
			assert.Equal(testcase.Expected, actual)
		})
	}
}
