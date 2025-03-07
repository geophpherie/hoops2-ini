package hoops2

import "slices"

type SortMethod int

const (
	SORT_TOTAL_POINTS SortMethod = iota
	SORT_NAME
	SORT_ROUND_1
	SORT_ROUND_2
	SORT_ROUND_3
	SORT_ROUND_4
	SORT_ROUND_5
	SORT_ROUND_6
	SORT_REGION_1
	SORT_REGION_2
	SORT_REGION_3
	SORT_REGION_4
)

func (p *Pool) SortResults(sortMethod SortMethod) {
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
	case SORT_ROUND_5:
		slices.SortFunc(p.Results, sortRound(5))
	case SORT_ROUND_6:
		slices.SortFunc(p.Results, sortRound(6))
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

	// find where TRUTH is
	var truth_index int
	for i, result := range p.Results {
		if result.Name == "TRUTH" {
			truth_index = i
			break
		}
	}

	// put it in front
	results := make([]PoolResult, 1, len(p.Results))
	results[0] = p.Results[truth_index]

	// fill in the rest
	results = append(results, p.Results[:truth_index]...)
	results = append(results, p.Results[truth_index+1:]...)

	p.Results = results
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
	// This sorts A -> Z
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
