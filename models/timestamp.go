package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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
	output := fmt.Sprint(time.Duration(*timestamp))
	output = strings.Replace(output, "h", ":", 1)
	output = strings.Replace(output, "m", ":", 1)
	output = strings.Replace(output, "s", "", 1)
	return []byte(output), nil
}
