package pacing

import (
	"encoding/json"
	"fmt"
	"time"
)

type conversionFunc func(float64) float64

// Unit is the representation of a rate of speed and contains functions to convert values to/from a standard unit (meters per sec)
type Unit struct {
	str          string
	duration     time.Duration
	distance     Distance
	toMPS        conversionFunc
	fromMPS      conversionFunc
	formatString func(float64) string
}

func (u Unit) String() string {
	return u.str
}

func (u Unit) Distance(val interface{}) (Distance, error) {
	switch v := val.(type) {
	case Distance:
		return v, nil
	case float64:
		return Distance(v) * u.distance, nil
	case string:
		return ParseDistance(v)
	default:
		return 0.0, fmt.Errorf("invalid distance type: %T", v)
	}

}

func (u Unit) DistanceInUnits(val interface{}) (float64, error) {
	d, err := u.Distance(val)
	if err != nil {
		return 0.0, err
	}
	switch u.distance {
	case Millimeter:
		return d.Millimeters(), nil
	case Centimeter:
		return d.Centimeters(), nil
	case Meter:
		return d.Meters(), nil
	case Kilometer:
		return d.Kilometers(), nil
	case Inch:
		return d.Inches(), nil
	case Yard:
		return d.Yards(), nil
	case Mile:
		return d.Miles(), nil
	default:
		return 0.0, fmt.Errorf("unsupported distance unit: %.2f", u.distance)
	}
}

func (u Unit) DistanceUnit() Distance {
	return u.distance
}

func (u Unit) DistanceString(val interface{}) string {
	d, err := u.DistanceInUnits(val)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%.2f %s", d, DistanceUnitString(u.distance))
}

func (u Unit) Duration(val interface{}) float64 {
	var d time.Duration
	switch v := val.(type) {
	case Duration:
		d = v.Duration
	case time.Duration:
		d = v
	default:
		panic(fmt.Sprintf("invalid duration type: %T", v))
	}

	switch u.duration {
	case time.Hour:
		return float64(d.Hours())
	case time.Minute:
		return float64(d.Minutes())
	case time.Second:
		return float64(d.Seconds())
	case time.Millisecond:
		return float64(d.Milliseconds())
	case time.Microsecond:
		return float64(d.Microseconds())
	}
	return float64(d)
}

func (u Unit) DurationUnit() string {
	switch u.duration {
	case time.Hour:
		return "hr"
	case time.Minute:
		return "min"
	case time.Second:
		return "sec"
	case time.Millisecond:
		return "ms"
	case time.Microsecond:
		return "μs"
	default:
		return "ns"
	}
}

func (u Unit) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *Unit) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		unit, err := Units(value)
		if err != nil {
			return err
		}
		*u = unit
		return nil
	default:
		return fmt.Errorf("invalid unit: %v. must use one of the following: %s", value, validUnits)
	}
}
