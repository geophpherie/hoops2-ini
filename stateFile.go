package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func loadStateData(filePath string) ([]string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(fileData), "\n")

	trimmedLines := make([]string, len(lines))

	for i, line := range lines {
		trimmedLines[i] = strings.TrimSpace(line)
	}

	return trimmedLines, nil
}

func parseStateData(fileLines []string) (Pool, error) {
	pool := Pool{}

	// first line is just a comment
	if fileLines[0] != "/* HOOPS2.INI state data file */" {
		return pool, errors.New("Invalid file header")
	}

	// second line has round weights, prefixed by 17 numbers (ignore for now)
	// last 6 are the round weights
	pool.RoundWeights = make(map[int]int)
	weights := strings.Fields(fileLines[1])
	for i, rw := range weights[len(weights)-6:] {
		round, err := strconv.Atoi(rw)
		if err != nil {
			panic(err)
		}
		pool.RoundWeights[i+1] = round
	}

	// lines 2,3,4,5 are regions
	pool.Regions = make(map[int]RegionDefinition)
	for i, regionLine := range fileLines[2:6] {
		region := RegionDefinition{}

		values := strings.Split(regionLine, "#")
		region.Name = values[0]

		region.Teams = make(map[int]string)
		for index, team := range values[1:] {
			if team == "" {
				continue
			}
			seed := index + 1
			region.Teams[seed] = team

		}

		pool.Regions[i+1] = region
	}

	// line 6 is the results, with the final 3 numbers referencing region number not seed
	// so region winner needs to be determined
	results := strings.Fields(fileLines[6])

	picks := []int{}
	for _, result := range results {
		seed, err := strconv.Atoi(result)
		if err != nil {
			panic(err)
		}
		picks = append(picks, seed)
	}
	pool.GameResults = NewPicks(picks)

	// the rest of the lines are contestant entries
	pool.Contestants = make(map[string]Contestant)
	for _, lineData := range fileLines[7:] {
		if strings.TrimSpace(lineData) == "" {
			continue
		}
		newContestant := Contestant{}
		fields := strings.Split(lineData, "#")

		newContestant.Name = fields[0]
		picks := []int{}
		for _, result := range strings.Fields(fields[1]) {
			seed, err := strconv.Atoi(result)
			if err != nil {
				panic(err)
			}
			picks = append(picks, seed)
		}
		newContestant.Picks = NewPicks(picks)

		pool.Contestants[newContestant.Name] = newContestant
	}

	return pool, nil
}

func exportStateData() {
	// TODO:
}
