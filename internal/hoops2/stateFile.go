package hoops2

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
		roundNumber := i + 1
		weight, err := strconv.Atoi(rw)
		if err != nil {
			panic(err)
		}
		pool.RoundWeights[roundNumber] = weight
	}

	// lines 2,3,4,5 are regions
	pool.Regions = make(map[int]Region)
	for i, regionLine := range fileLines[2:6] {
		regionNumber := i + 1
		region := Region{}

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

		pool.Regions[regionNumber] = region
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

func ExportToJson(pool Pool, path string) error {
	data, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(path)
	// Create the directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
		fmt.Println("Directory created successfully:", dirPath)
	} else if err != nil {
		fmt.Println("Error checking directory:", err)
		return err
	}

	err = os.WriteFile(path, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
