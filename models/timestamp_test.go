package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(test *testing.T) {
	assert := assert.New(test)

	numbers := []struct {
		Input    string
		Expected TimeStamp
	}{
		{Input: "1500", Expected: 1500},
		{Input: "3.54", Expected: 3},
		{Input: "600", Expected: 600},
		// Should we allow negative numbers?
	}
	for _, testcase := range numbers {
		test.Run(
			fmt.Sprintf("Should convert to %v any numeric value without error from %s", testcase.Expected, testcase.Input),
			func(test *testing.T) {
				// Arrange
				bytes := []byte(testcase.Input)
				actual := TimeStamp(0)
				timestamp := &actual

				// Act
				exception := timestamp.UnmarshalJSON(bytes)

				// Assert
				assert.Nil(exception)
				assert.Equal(testcase.Expected, actual)
			},
		)
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
		test.Run(
			fmt.Sprintf("Should convert a valid quoted string value to %v without error from %s", testcase.Expected, testcase.Input),
			func(test *testing.T) {
				// Arrange
				bytes := []byte(testcase.Input)
				actual := TimeStamp(0)
				timestamp := &actual

				// Act
				exception := timestamp.UnmarshalJSON(bytes)

				// Assert
				assert.Nil(exception)
				assert.Equal(testcase.Expected, actual)
			},
		)
	}

	invalids := []struct {
		Input  string
		Reason string
	}{
		// Invalid JSON values (JSON Parser Error)
		{Input: `"1500`, Reason: "unexpected end of JSON input"},
		{Input: `"3h4m5s`, Reason: "unexpected end of JSON input"},
		{Input: `0900"`, Reason: "invalid character"},
		{Input: "15:00", Reason: "invalid character"},
		{Input: "10:45:15", Reason: "invalid character"},
		{Input: "3h4m5s", Reason: "invalid character"},
		{Input: "1:2:3", Reason: "invalid character"},
		{Input: `1:2:3"`, Reason: "invalid character"},
		{Input: "'3h4m5s'", Reason: "invalid character"},
		{Input: `["1500"]`, Reason: "invalid time stamp"},

		// Valid JSON string value, but not a time.Duration
		{Input: `"15,00"`, Reason: `unknown unit`},
		{Input: `"hello-world"`, Reason: "time: invalid duration"},
		{Input: `"(3h4m5s)"`, Reason: "time: invalid duration"},

		// Valid JSON value, but not float nor string
		{Input: `{"1500": "0600"}`, Reason: "invalid time stamp"},
		{Input: "[1500]", Reason: "invalid time stamp"},
		{Input: `{"1500": 600}`, Reason: "invalid time stamp"},
		{Input: "true", Reason: "invalid time stamp"},
		{Input: "false", Reason: "invalid time stamp"},
	}
	for _, testcase := range invalids {
		test.Run(
			fmt.Sprintf("Should NOT convert an invalid string value (%v) returning error (%s)", testcase.Input, testcase.Reason),
			func(test *testing.T) {
				// Arrange
				bytes := []byte(testcase.Input)
				actual := TimeStamp(0)
				timestamp := &actual

				// Act
				exception := timestamp.UnmarshalJSON(bytes)

				// Assert
				assert.NotNil(exception)
				assert.Contains(exception.Error(), testcase.Reason)
			},
		)
	}
}

func TestMarshalJSON(test *testing.T) {
	assert := assert.New(test)
	testcases := []struct {
		Input    TimeStamp
		Expected string
	}{
		{Input: 600, Expected: `"00:10:00"`},
		{Input: 3600, Expected: `"01:00:00"`},
		{Input: 45, Expected: `"00:00:45"`},
		{Input: 7, Expected: `"00:00:07"`},
		{Input: 36359, Expected: `"10:05:59"`},
		{Input: 24*3600 + 5*60 + 1, Expected: `"24:05:01"`},
		{Input: 132*3600 + 5*60 + 1, Expected: `"132:05:01"`},
	}
	for _, testcase := range testcases {
		test.Run(fmt.Sprintf("Should format correctly from %v to '%s'", testcase.Input, testcase.Expected), func(test *testing.T) {
			// Arrange
			time := &testcase.Input
			expected := []byte(testcase.Expected)

			// Act
			actual, exception := time.MarshalJSON()

			// Assert
			assert.Nil(exception)
			assert.Equal(expected, actual)
		})
	}
}
