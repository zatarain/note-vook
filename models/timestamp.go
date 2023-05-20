package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type TimeStamp int64

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
		pattern := regexp.MustCompile("^(([0-9]+):)?([0-9]+):([0-9]+)$")
		output := pattern.ReplaceAllString(value, "0${2}h${3}m${4}s")
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
	zero, _ := time.Parse(time.TimeOnly, "00:00:00")
	output := fmt.Sprintf("\"%v\"", zero.Add(duration).Format(time.TimeOnly))
	return []byte(output), nil
}
