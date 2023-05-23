package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TimeStamp int64

const (
	PATTERN     string = "^(([0-9]+):)?([0-9]+):([0-9]+)$"
	REPLACEMENT string = "0${2}h${3}m${4}s"
	ZERO        string = "00:00:00"
)

func (timestamp *TimeStamp) UnmarshalJSON(bytes []byte) error {
	var unmarshalledJson interface{}

	exception := json.Unmarshal(bytes, &unmarshalledJson)
	if exception != nil {
		return exception
	}

	switch value := unmarshalledJson.(type) {
	case float64:
		*timestamp = TimeStamp(value)
	case string:
		number, cannotConvert := strconv.ParseFloat(value, 64)
		if cannotConvert == nil {
			*timestamp = TimeStamp(number)
			break
		}

		pattern := regexp.MustCompile(PATTERN)
		output := pattern.ReplaceAllString(value, REPLACEMENT)
		duration, exception := time.ParseDuration(output)
		if exception != nil {
			return exception
		}
		*timestamp = TimeStamp(duration / time.Second)
	default:
		return fmt.Errorf("invalid time stamp: %#v", unmarshalledJson)
	}

	return nil
}

func (timestamp *TimeStamp) MarshalJSON() ([]byte, error) {
	value := int64(*timestamp)
	duration := time.Duration(value) * time.Second
	zero, _ := time.Parse(time.TimeOnly, ZERO)
	output := zero.Add(duration).Format(time.TimeOnly)
	if duration >= 24*time.Hour {
		hours := int64(duration / time.Hour)
		output = regexp.MustCompile("^([0-9]{2})").ReplaceAllString(output, fmt.Sprintf("%d", hours))
	}
	return []byte(fmt.Sprintf(`"%v"`, output)), nil
}
