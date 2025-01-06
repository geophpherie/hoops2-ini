package main

import (
	"fmt"
	"path/filepath"
	"slices"
)

type Pool struct {
	Name         string
	RoundWeights map[int]int
	Regions      map[int]RegionDefinition
	GameResults  Picks
	Contestants  map[string]Contestant
	Results      []PoolResult
}

type RegionDefinition struct {
	Name  string
	Teams map[int]string
}

type Contestant struct {
	Name  string
	Picks Picks
}

type PoolResult struct {
	Name string
	PoolScore
	FinishingPicks
}

const (
	SORT_TOTAL_POINTS = iota
	SORT_NAME
	SORT_ROUND_1
	SORT_ROUND_2
	SORT_ROUND_3
	SORT_ROUND_4
	SORT_REGION_1
	SORT_REGION_2
	SORT_REGION_3
	SORT_REGION_4
)

func (p *Pool) SortResults(sortMethod int) {
	switch sortMethod {
	case SORT_TOTAL_POINTS:
		slices.SortFunc(p.Results, sortTotalPoints)
	case SORT_NAME:
		slices.SortFunc(p.Results, sortName)
	case SORT_ROUND_1:
		slices.SortFunc(p.Results, sortRound(1))
	case SORT_ROUND_2:
		slices.SortFunc(p.Results, sortRound(2))
	case SORT_ROUND_3:
		slices.SortFunc(p.Results, sortRound(3))
	case SORT_ROUND_4:
		slices.SortFunc(p.Results, sortRound(4))
	case SORT_REGION_1:
		slices.SortFunc(p.Results, sortRegion(1))
	case SORT_REGION_2:
		slices.SortFunc(p.Results, sortRegion(2))
	case SORT_REGION_3:
		slices.SortFunc(p.Results, sortRegion(3))
	case SORT_REGION_4:
		slices.SortFunc(p.Results, sortRegion(4))
	default:
		panic("unknown sort")
	}
}

func (pr PoolResult) Format() string {
	totalScore := fmt.Sprintf("%3v", pr.Total)
	name := fmt.Sprintf("%-6v", pr.Name)

	round1 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[1], pr.CumulativeScoreByRound[1])
	round2 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[2], pr.CumulativeScoreByRound[2])
	round3 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[3], pr.CumulativeScoreByRound[3])
	round4 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[4], pr.CumulativeScoreByRound[4])
	round5 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[5], pr.CumulativeScoreByRound[5])
	round6 := fmt.Sprintf("%3v/%3v", pr.CumulativeCorrectPicksByRound[6], pr.CumulativeScoreByRound[6])

	region1 := fmt.Sprintf(
		"%3v/%3v/%3v/%3v",
		pr.CumulativeScoreByRoundByRegion[1][1],
		pr.CumulativeScoreByRoundByRegion[1][2],
		pr.CumulativeScoreByRoundByRegion[1][3],
		pr.CumulativeScoreByRoundByRegion[1][4],
	)

	region2 := fmt.Sprintf(
		"%3v/%3v/%3v/%3v",
		pr.CumulativeScoreByRoundByRegion[2][1],
		pr.CumulativeScoreByRoundByRegion[2][2],
		pr.CumulativeScoreByRoundByRegion[2][3],
		pr.CumulativeScoreByRoundByRegion[2][4],
	)

	region3 := fmt.Sprintf(
		"%3v/%3v/%3v/%3v",
		pr.CumulativeScoreByRoundByRegion[3][1],
		pr.CumulativeScoreByRoundByRegion[3][2],
		pr.CumulativeScoreByRoundByRegion[3][3],
		pr.CumulativeScoreByRoundByRegion[3][4],
	)

	region4 := fmt.Sprintf(
		"%3v/%3v/%3v/%3v",
		pr.CumulativeScoreByRoundByRegion[4][1],
		pr.CumulativeScoreByRoundByRegion[4][2],
		pr.CumulativeScoreByRoundByRegion[4][3],
		pr.CumulativeScoreByRoundByRegion[4][4],
	)

	regionFinalists := fmt.Sprintf(
		"%2v %2v %2v %2v %2v %2v %2v %2v",
		pr.RegionFinalists[0],
		pr.RegionFinalists[1],
		pr.RegionFinalists[2],
		pr.RegionFinalists[3],
		pr.RegionFinalists[4],
		pr.RegionFinalists[5],
		pr.RegionFinalists[6],
		pr.RegionFinalists[7],
	)

	finalFour := fmt.Sprintf(
		"%2v %2v %2v %2v ",
		pr.FinalFour[0],
		pr.FinalFour[1],
		pr.FinalFour[2],
		pr.FinalFour[3],
	)

	finals := fmt.Sprintf("%v%v%v", pr.FinalGame[0], pr.FinalGame[1], pr.Winner)

	return fmt.Sprintf(
		"%v %v %v %v %v %v %v %v %v %v %v %v  %v  %v %v",
		totalScore,
		name,
		round1,
		round2,
		round3,
		round4,
		round5,
		round6,
		region1,
		region2,
		region3,
		region4,
		regionFinalists,
		finalFour,
		finals,
	)
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

	pool.Name = filepath.Ext(filePath)[2:]

	score := CalculateScoring(pool.GameResults, pool.RoundWeights)
	poolResults := []PoolResult{
		{
			Name:           "TRUTH",
			PoolScore:      score,
			FinishingPicks: pool.GameResults.FinishingPicks(),
		},
	}

	for name, contestant := range pool.Contestants {
		comparedPicks := ComparePicks(pool.GameResults, contestant.Picks)
		score = CalculateScoring(comparedPicks, pool.RoundWeights)

		poolResults = append(poolResults, PoolResult{
			Name:           name,
			PoolScore:      score,
			FinishingPicks: contestant.Picks.FinishingPicks(),
		})
	}

	pool.Results = poolResults

	return pool, nil
}
