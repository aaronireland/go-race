package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/aaronireland/go-race/grid"
	"github.com/aaronireland/go-race/pacing"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var gridFilePath *string

func init() {
	log.SetHandler(text.New(os.Stderr))
	gridFilePath = flag.String("grid", "grid.json", "Path to json file that contains the requested pace segments")
}

func readGridFromFile(filePath string) (grid.CourseGrid, error) {
	var course grid.CourseGrid
	gridFile, err := os.Open(filePath)
	if err != nil {
		return course, err
	}
	defer func() { _ = gridFile.Close() }()
	gridFileContent, err := ioutil.ReadAll(gridFile)
	if err != nil {
		return course, err
	}

	err = json.Unmarshal(gridFileContent, &course)
	return course, err
}

func finishRace(race grid.Race, totalDistance pacing.Distance, maxTime pacing.Duration) grid.CourseSegment {
	var (
		distanceCompleted pacing.Distance
		timeElapsed       pacing.Duration
	)

	for _, segment := range race.Segments {
		_, distance, duration := segment.Stats()
		distanceCompleted += distance
		timeElapsed = timeElapsed.Add(duration)
	}

	timeRemaining := maxTime.Subtract(timeElapsed)
	distanceRemaining := totalDistance - distanceCompleted

	pace := pacing.Calculate(timeRemaining, distanceRemaining, race.Units)

	return grid.CourseSegment{
		Pace:        pace,
		Distance:    distanceRemaining,
		TimeElapsed: timeRemaining,
		Units:       race.Units,
	}
}

func main() {
	flag.Parse()
	grid, err := readGridFromFile(*gridFilePath)
	if err != nil {
		log.WithError(err).Fatal("Invalid pace grid file")
	}
	courseDistance, err := pacing.ParseDistance(grid.TotalDistance)
	if err != nil {
		log.WithError(err).Fatal("Invalid pace grid file")
	}

	races := "races"
	if len(grid.Races) == 1 {
		races = "race"
	}

	log.WithField("totalDistance", grid.TotalDistance).WithField("totalDuration", grid.TotalDuration).Infof("Parsed pacing grid with %d %s...", len(grid.Races), races)

	for i, race := range grid.Races {
		remainingSegment := finishRace(race, courseDistance, grid.TotalDuration)
		for si, segment := range race.Segments {
			pace, distance, duration := segment.Stats()
			d := race.Units.DistanceString(distance)
			log.WithField("race", i+1).WithField("segment", si+1).Infof("%s completed in %s with a pace of %s %s", d, duration, pace, pace.Units)
		}
		pace, distance, duration := remainingSegment.Stats()
		d := race.Units.DistanceString(distance)
		log.WithField("race", i+1).Infof(
			"A pace of %s %s is required to complete the remaining %s of %s "+
				"in %s to achieve a time of %s",
			pace, pace.Units, d, grid.TotalDistance,
			duration, grid.TotalDuration,
		)
	}
}
