package pacing

import (
	"fmt"
	"strconv"
	"time"
)

// Pace contextualizes a pacing unit with a value and provides methods for extracting distances, times, and concatenation
type Pace struct {
	Rate  float64
	Units Unit
}

// New is the constructor for a Pace. It attempts to parses the pace rate as a duration and converts
// that duration to a decimal value using the units provided
func New(pace interface{}, units Unit) (Pace, error) {
	switch val := pace.(type) {
	case float64:
		return Pace{val, units}, nil
	case Duration:
		return Pace{units.Duration(val), units}, nil
	case time.Duration:
		return Pace{units.Duration(Duration{val}), units}, nil
	case string:
		duration, err := ParseDuration(val)
		if err != nil {
			if paceAsFloat, err := strconv.ParseFloat(val, 64); err == nil {
				return Pace{paceAsFloat, units}, nil
			}
			return Pace{}, fmt.Errorf("invalid pace rate %s: %w", val, err)
		}
		return Pace{units.Duration(duration), units}, nil
	default:
		return Pace{}, fmt.Errorf("unsupported pace value type: %T", val)
	}
}

// Calculate extrapolates the pace value from a duration and distance value (distance units gathered from the Units provided)
func Calculate(duration Duration, distance Distance, units Unit) Pace {
	return Pace{units.fromMPS(distance.Meters() / float64(duration.Seconds())), units}
}

func (p Pace) String() string {
	return p.Units.formatString(p.Rate)
}

func (p Pace) mps() float64 {
	return p.Units.toMPS(p.Rate)
}

func (p Pace) Distance(t Duration) Distance {
	metersTraveled := p.mps() * float64(t.Duration.Seconds())
	return Distance(metersTraveled) * Meter
}

func (p Pace) Duration(distance Distance) Duration {
	secondsElapsed := distance.Meters() / p.mps()
	return Duration{time.Duration(secondsElapsed) * time.Second}
}

func (p Pace) As(units Unit) Pace {
	return Pace{units.fromMPS(p.mps()), units}
}
