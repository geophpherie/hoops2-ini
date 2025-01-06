package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Need to pass path to INIT file")
		os.Exit(1)
	}

	pool, err := NewPoolFromFile(args[0])
	if err != nil {
		panic(err)
	}

	pool.SortResults(SORT_NAME)
	for _, v := range pool.Results {
		fmt.Println(v.Format())
	}

	pool.SortResults(SORT_TOTAL_POINTS)
	for _, v := range pool.Results {
		fmt.Println(v.Format())
	}
}

func sortTotalPoints(a, b PoolResult) int {
	// this sorts descending numerical
	if a.Total < b.Total {
		return 1
	} else if a.Total > b.Total {
		return -1
	} else {
		return 0
	}
}

func sortName(a, b PoolResult) int {
	// This sorts ! -> Z
	if a.Name < b.Name {
		return -1
	} else if a.Name > b.Name {
		return 1
	} else {
		return 0
	}
}

func sortRound(roundNumber int) func(a, b PoolResult) int {
	return func(a, b PoolResult) int {
		// this sorts descending numerical
		if a.CumulativeScoreByRound[roundNumber] < b.CumulativeScoreByRound[roundNumber] {
			return 1
		} else if a.CumulativeScoreByRound[roundNumber] > b.CumulativeScoreByRound[roundNumber] {
			return -1
		} else {
			return 0
		}
	}
}

func sortRegion(regionNumber int) func(a, b PoolResult) int {
	return func(a, b PoolResult) int {
		// this sorts descending numerical
		if a.CumulativeScoreByRoundByRegion[regionNumber][4] < b.CumulativeScoreByRoundByRegion[regionNumber][4] {
			return 1
		} else if a.CumulativeScoreByRoundByRegion[regionNumber][4] > b.CumulativeScoreByRoundByRegion[regionNumber][4] {
			return -1
		} else {
			return 0
		}
	}
}
