package hoops2

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Pool struct {
	Name         string
	RoundWeights map[int]int
	Regions      map[int]Region
	GameResults  Picks
	Contestants  map[string]Contestant
	Results      []PoolResult
	LastModified string
}

type Region struct {
	Name  string
	Teams map[int]string
}

type Contestant struct {
	Name  string
	Picks Picks
}

func NewPoolFromFile(filePath string) (Pool, error) {
	fileLines, err := loadStateData(filePath)
	if err != nil {
		return Pool{}, err
	}

	pool, err := parseStateData(fileLines)
	if err != nil {
		return Pool{}, err
	}

	year, err := strconv.Atoi(filepath.Ext(filePath)[2:])
	if err != nil {
		return Pool{}, err
	}
	if year < 85 {
		pool.Name = fmt.Sprintf("20%02d", year)
	} else {
		pool.Name = fmt.Sprintf("19%02d", year)
	}

	score := CalculateScoring(pool.GameResults, pool.RoundWeights)

	// use TRUTH as the first result in the Pool
	poolResults := []PoolResult{
		{
			Name:           "TRUTH",
			PoolScore:      score,
			FinishingPicks: pool.GameResults.FinishingPicks(),
		},
	}

	// add the rest of the scoring results comparing to TRUTH
	for name, contestant := range pool.Contestants {
		comparedPicks := pool.GameResults.Compare(contestant.Picks)
		score = CalculateScoring(comparedPicks, pool.RoundWeights)

		poolResults = append(poolResults, PoolResult{
			Name:           name,
			PoolScore:      score,
			FinishingPicks: contestant.Picks.FinishingPicks(),
		})
	}

	pool.Results = poolResults

	fileInfo, _ := os.Stat(filePath)
	pool.LastModified = fileInfo.ModTime().Format("January 2, 2006 at 03:04:05 PM MST")

	return pool, nil
}
