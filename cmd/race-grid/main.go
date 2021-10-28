package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/aaronireland/go-race/grid"
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

func main() {
	flag.Parse()
	grid, err := readGridFromFile(*gridFilePath)
	if err != nil {
		log.WithError(err).Fatal("Invalid pace grid file")
	}
	log.WithField("totalDistance", grid.TotalDistance).WithField("totalDuration", grid.TotalDuration).Infof("Parsed pacing grid with %d segments...", len(grid.Races))

	/*
		var (
			distanceRemaining float64 = grid.TotalDistance
			timeRemaining pacing.Duration = grid.TotalTime
		)

	*/
}
