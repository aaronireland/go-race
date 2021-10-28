package grid

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aaronireland/go-race/pacing"
)

type CourseSegment struct {
	Pace        pacing.Pace     `json:"-"`
	Distance    pacing.Distance `json:"-"`
	TimeElapsed pacing.Duration `json:"time,omitempty"`
	Units       pacing.Unit     `json:"units"`
}

type courseSegmentJSON struct {
	Pace        string          `json:"pace"`
	Distance    float64         `json:"distance,omitempty"`
	TimeElapsed pacing.Duration `json:"time,omitempty"`
	Units       pacing.Unit     `json:"units"`
}

type Race struct {
	Segments []CourseSegment `json:"segments"`
	Units    pacing.Unit     `json:"units"`
}

type CourseGrid struct {
	TotalDistance string          `json:"totalDistance"`
	TotalDuration pacing.Duration `json:"totalTime"`
	Races         []Race          `json:"races"`
}

type courseGridJSON CourseGrid

func (g *CourseGrid) UnmarshalJSON(data []byte) error {
	var (
		gridData courseGridJSON
		grid     CourseGrid
	)

	if err := json.Unmarshal(data, &gridData); err != nil {
		return err
	}

	grid = CourseGrid(gridData)
	courseDistance, err := pacing.ParseDistance(grid.TotalDistance)
	if err != nil {
		return fmt.Errorf("invalid distance for course grid: %w", err)
	}

	for _, race := range grid.Races {
		var (
			totalDistance pacing.Distance
			totalDuration time.Duration
		)
		for _, segment := range race.Segments {
			if segment.Distance > 0 {
				totalDistance += segment.Distance
				totalDuration += segment.Pace.Duration(segment.Distance).Time()
			} else {
				totalDistance += segment.Pace.Distance(segment.TimeElapsed)
				totalDuration += segment.TimeElapsed.Time()
			}

		}

		if totalDistance >= courseDistance || totalDuration >= grid.TotalDuration.Time() {
			dist, _ := race.Units.DistanceInUnits(totalDistance)
			courseDist, _ := race.Units.DistanceInUnits(courseDistance)
			unit := pacing.DistanceUnitString(race.Units.DistanceUnit())
			return fmt.Errorf("invalid race in grid: course segments total %.2f%s and %s where course totals %.2f%s and %s", dist, unit, totalDuration, courseDist, unit, grid.TotalDuration)
		}
	}

	*g = grid
	return nil
}

func (csj courseSegmentJSON) toCourseSegment() (CourseSegment, error) {
	var s = CourseSegment{
		TimeElapsed: csj.TimeElapsed,
		Units:       csj.Units,
	}

	s.Distance = pacing.Distance(csj.Distance) * s.Units.DistanceUnit()
	pace, err := csj.calculatePace()
	if err != nil {
		return s, err
	}
	s.Pace = pace
	return s, nil
}

func (s CourseSegment) toJSON() (courseSegmentJSON, error) {
	distance, err := s.Units.DistanceInUnits(s.Distance)
	if err != nil {
		return courseSegmentJSON{}, err
	}

	return courseSegmentJSON{
		Pace:        s.Pace.String(),
		Distance:    distance,
		TimeElapsed: s.TimeElapsed,
		Units:       s.Units,
	}, nil
}

func (s CourseSegment) MarshalJSON() ([]byte, error) {
	data, err := s.toJSON()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(data)
}

func (s *CourseSegment) UnmarshalJSON(data []byte) error {
	var segmentJSON courseSegmentJSON

	if err := json.Unmarshal(data, &segmentJSON); err != nil {
		return err
	}

	segment, err := segmentJSON.toCourseSegment()
	if err != nil {
		return err
	}
	*s = segment
	return nil
}

func (csj courseSegmentJSON) calculatePace() (pacing.Pace, error) {
	var pace pacing.Pace

	// Validate that enough pacing data was provided for the segment
	okPace := (csj.Pace != "" && (csj.Distance > 0 || csj.TimeElapsed.Duration > 0))
	okDistanceDuration := (csj.Pace == "" && csj.Distance > 0 && csj.TimeElapsed.Duration > 0)
	distanceStr := fmt.Sprintf("%.2f %s", csj.Distance, pacing.DistanceUnitString(csj.Units.DistanceUnit()))
	if !(okPace || okDistanceDuration) {
		return pace, fmt.Errorf("invalid course segment: (pace=%s), (distance=%s), (duration=%s)", csj.Pace, distanceStr, csj.TimeElapsed)
	}

	if csj.Pace != "" {
		return pacing.Parse(csj.Pace, csj.Units)
	}

	distance, err := csj.Units.Distance(csj.Distance)
	if err != nil {
		return pace, err
	}
	return pacing.Calculate(csj.TimeElapsed, distance, csj.Units), nil
}
