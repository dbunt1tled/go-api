package timestamp

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var str string
	if err := sonic.ConfigFastest.Unmarshal(data, &str); err == nil {
		parsedTime, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
		t.Time = parsedTime
		return nil
	}

	var unixTime int64
	if err := sonic.ConfigFastest.Unmarshal(data, &unixTime); err == nil {
		t.Time = time.Unix(unixTime, 0)
		return nil
	}

	return errors.New("timestamp: invalid format")
}

func (t *Timestamp) Scan(value any) error {
	if _, ok := value.(time.Time); ok {
		t.Time = value.(time.Time)
		return nil
	}

	return errors.New("timestamp: invalid format")
}
