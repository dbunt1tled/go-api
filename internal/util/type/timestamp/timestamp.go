package timestamp

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		parsedTime, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
		t.Time = parsedTime
		return nil
	}

	var unixTime int64
	if err := json.Unmarshal(data, &unixTime); err == nil {
		t.Time = time.Unix(unixTime, 0)
		return nil
	}

	return errors.New("timestamp: invalid format")
}
