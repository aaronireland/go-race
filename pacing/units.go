package pacing

import (
	"fmt"
	"strings"
	"time"
)

const (
	metersPerKm   = 1000.00
	secondsPerHr  = 3600.00
	minutesPerHr  = 60.00
	secondsPerMin = 60.00
	metersPerMile = 1609.34
	minMilePerMps = 26.8224
	minKmPerMps   = float64(50.0 / 3.0)
	kphPerMps     = 3.6
	mphPerMps     = 2.2369362920544
)

type units map[string]Unit

func (u units) String() string {
	unitNames := make([]string, 0, len(u))
	for unit := range u {
		unitNames = append(unitNames, unit)
	}
	return strings.Join(unitNames, ", ")
}

var validUnits = units{
	"kph":      KPH,
	"mph":      MPH,
	"min/km":   MINKM,
	"min/mile": MINMILE,
}

var fmtFloat64 = func(rate float64) string {
	return fmt.Sprintf("%.2f", rate)
}

func fmtTimeDuration(units time.Duration) func(float64) string {
	return func(rate float64) string {
		var str string
		d := time.Duration(rate * float64(units))
		hours := int(d.Hours())
		d -= (time.Duration(hours) * time.Hour)

		minutes := int(d.Minutes())
		d -= (time.Duration(minutes) * time.Minute)

		seconds := int(d.Seconds())

		if hours == 0 {
			str += fmt.Sprintf("%d:%02d", minutes, seconds)
		} else {
			str += fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
		}
		return str
	}

}

func Units(u string) (Unit, error) {
	if unit, ok := validUnits[u]; ok {
		return unit, nil
	}
	return Unit{}, fmt.Errorf("invalid unit: %s. must use one of the following: %s", u, validUnits)
}

// KPH represents rates/distance as kilometers per hour
var KPH = Unit{
	"kph",
	time.Hour,
	Kilometer,
	func(kph float64) float64 {
		return kph * metersPerKm / secondsPerHr
	},
	func(mps float64) float64 {
		return mps * kphPerMps
	},
	fmtFloat64,
}

// MPH represents rates/distance as miles per hour
var MPH = Unit{
	"mph",
	time.Hour,
	Mile,
	func(mph float64) float64 {
		return mph * metersPerMile / secondsPerHr
	},
	func(mps float64) float64 {
		return mps * mphPerMps
	},
	fmtFloat64,
}

// MINKM represents rates/distance as minutes per kilometer
var MINKM = Unit{
	"min/km",
	time.Minute,
	Kilometer,
	func(minkm float64) float64 {
		return minKmPerMps / minkm
	},
	func(mps float64) float64 {
		return minKmPerMps / mps
	},
	fmtTimeDuration(time.Minute),
}

// MINMILE represents rates/distances as minutes per mile
var MINMILE = Unit{
	"min/mile",
	time.Minute,
	Mile,
	func(minmile float64) float64 {
		return minMilePerMps / minmile
	},
	func(mps float64) float64 {
		return minMilePerMps / mps
	},
	fmtTimeDuration(time.Minute),
}
