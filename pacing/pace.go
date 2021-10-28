package pacing

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Pace contextualizes a pacing unit with a value and provides methods for extracting distances, times, and concatenation
type Pace struct {
	Rate  float64
	Units Unit
}

// New is the basic constructor for a Pace
func New(pace float64, units Unit) Pace {
	return Pace{pace, units}
}

// Calculate extrapolates the pace value from a duration and distance value (distance units gathered from the Units provided)
func Calculate(duration Duration, distance Distance, units Unit) Pace {
	return Pace{units.fromMPS(distance.Meters() / float64(duration.Seconds())), units}
}

// Parse attempts to convert a string value into a pace value (either a time duration or a decimal distance). Valid
// inputs are as such: hh:mm:ss, mm:ss, 00h00m00s (or any valid time.Duration string), or a numeric value like 1.23 etc.
func Parse(pace string, units Unit) (Pace, error) {
	parts := strings.Split(pace, ":")
	switch len(parts) {
	case 3: // hh:mm:ss
		pace = fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2])
		return parseFromDuration(pace, units)
	case 2: // mm:ss
		pace = fmt.Sprintf("%sm%ss", parts[0], parts[1])
		return parseFromDuration(pace, units)
	default: // either 00h00m00s, 00m00s, 00s, or a decimal value
		p, err := parseFromDuration(pace, units)
		if err != nil {
			if paceAsFloat, err := strconv.ParseFloat(pace, 64); err == nil {
				return New(paceAsFloat, units), nil
			}
			return Pace{}, fmt.Errorf("invalid pace string \"%s\": %w", pace, err)
		}
		return p, nil
	}
}

// MustParse attempts to convert a string value in a pace value and panics if the input is invalid (see: pacing.Parse)
func MustParse(pace string, units Unit) Pace {
	p, err := Parse(pace, units)
	if err != nil {
		panic(err)
	}
	return p
}

func parseFromDuration(pace string, units Unit) (Pace, error) {
	paceTime, err := time.ParseDuration(pace)
	if err != nil {
		return Pace{}, fmt.Errorf("invalid pace string: %w", err)
	}
	return New(units.Duration(paceTime), units), nil
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
