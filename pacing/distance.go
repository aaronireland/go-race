package pacing

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	Millimeter Distance = 1.0
	Centimeter          = 10 * Millimeter
	Meter               = 100 * Centimeter
	Kilometer           = 1000 * Meter
	Inch                = 25.4 * Millimeter
	Yard                = 36 * Inch
	Mile                = 1760 * Yard
)

type Distance float64

func DistanceUnitString(d Distance) string {
	switch d {
	case Millimeter:
		return "mm"
	case Centimeter:
		return "cm"
	case Meter:
		return "m"
	case Kilometer:
		return "km"
	case Inch:
		return "in"
	case Yard:
		return "yd"
	case Mile:
		return "mi"
	default:
		return ""
	}
}

func ParseDistance(distStr string) (Distance, error) {
	distanceRe := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	unitsRe := regexp.MustCompile(`mi|km|yd|in|cm|mm|m`)

	if units := unitsRe.FindString(distStr); units != "" {
		distance, err := strconv.ParseFloat(distanceRe.FindString(distStr), 64)
		if err != nil {
			return 0.0, err
		}
		d := Distance(distance)
		switch units {
		case "mi":
			return d * Mile, nil
		case "yd":
			return d * Yard, nil
		case "in":
			return d * Inch, nil
		case "km":
			return d * Kilometer, nil
		case "m":
			return d * Meter, nil
		case "cm":
			return d * Centimeter, nil
		case "mm":
			return d * Millimeter, nil
		default:
			return 0.0, fmt.Errorf("invalid distance units: %s", units)

		}
	}
	return 0.0, fmt.Errorf("invalid distance: %s", distStr)
}

func MustParseDistance(distStr string) Distance {
	dist, err := ParseDistance(distStr)
	if err != nil {
		panic(err)
	}
	return dist
}

func (d Distance) Millimeters() float64 {
	return float64(d)
}

func (d Distance) Centimeters() float64 {
	return float64(d / Centimeter)
}

func (d Distance) Meters() float64 {
	return float64(d / Meter)
}

func (d Distance) Kilometers() float64 {
	return float64(d / Kilometer)
}

func (d Distance) Inches() float64 {
	return float64(d / Inch)
}

func (d Distance) Yards() float64 {
	return float64(d / Yard)
}

func (d Distance) Miles() float64 {
	return float64(d / Mile)
}
