package pacing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

// ParseDuration attempts to convert a string value into a time duration Valid inputs are something like:
// hh:mm:ss, mm:ss, 00h00m00s (or any valid time.Duration string).
func ParseDuration(d string) (Duration, error) {
	parts := strings.Split(d, ":")
	switch len(parts) {
	case 3: // hh:mm:ss
		d = fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2])
	case 2: // mm:ss
		d = fmt.Sprintf("%sm%ss", parts[0], parts[1])
	}
	duration, err := time.ParseDuration(d)
	return Duration{duration}, err
}

func (d Duration) Add(duration Duration) Duration {
	return Duration{d.Duration + duration.Duration}
}

func (d Duration) Subtract(duration Duration) Duration {
	remaining := d.Time() - duration.Time()
	if remaining < 0 {
		return Duration{0}
	}
	return Duration{remaining}
}

func (d Duration) Time() time.Duration {
	return d.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		*d, err = ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
