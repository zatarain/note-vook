package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimeStamp struct {
	time.Duration
}

func (timestamp *TimeStamp) UnmarshalJSON(b []byte) error {
	var unmarshalledJson interface{}

	exception := json.Unmarshal(b, &unmarshalledJson)
	if exception != nil {
		return exception
	}

	switch value := unmarshalledJson.(type) {
	case float64:
		timestamp.Duration = time.Duration(value)
	case string:
		timestamp.Duration, exception = time.ParseDuration(value)
		if exception != nil {
			return exception
		}
	default:
		return fmt.Errorf("invalid time stamp: %#v", unmarshalledJson)
	}

	return nil
}
