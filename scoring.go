package main

type PoolScore struct {
	Total                          int
	CumulativeCorrectPicksByRound  map[int]int
	CumulativeScoreByRound         map[int]int
	CumulativeScoreByRoundByRegion map[int]map[int]int
}

type RoundScore struct {
	Games int
	Score int
}

func (rs RoundScore) Add(other RoundScore) RoundScore {
	return RoundScore{
		Games: rs.Games + other.Games,
		Score: rs.Score + other.Score,
	}
}

func CalculateScoring(picks Picks, roundWeights map[int]int) PoolScore {
	cumCorrectPicks := make(map[int]int)
	cumRoundScore := make(map[int]int)

	cumScore := RoundScore{Games: 0, Score: 0}
	for i := 1; i <= 6; i++ {
		roundScore := CalculateRoundScore(picks.RoundSeeds(i), roundWeights[i])
		cumScore = roundScore.Add(cumScore)

		cumCorrectPicks[i] = cumScore.Games
		cumRoundScore[i] = cumScore.Score
	}

	cumRegionScore := make(map[int]map[int]int)
	// all regions
	for i := 1; i <= 4; i++ {
		cumRegionScore[i] = map[int]int{}
		cumScore := RoundScore{Games: 0, Score: 0}

		// all rounds in region
		for j := 1; j <= 4; j++ {
			roundScore := CalculateRoundScore(picks.Region(i)[j], roundWeights[j])
			cumScore = roundScore.Add(cumScore)
			cumRegionScore[i][j] = cumScore.Score
		}
	}

	return PoolScore{
		Total:                          cumRoundScore[6],
		CumulativeCorrectPicksByRound:  cumCorrectPicks,
		CumulativeScoreByRound:         cumRoundScore,
		CumulativeScoreByRoundByRegion: cumRegionScore,
	}
}

func CalculateRoundScore(winningSeeds []int, roundWeight int) RoundScore {
	roundScore := RoundScore{Games: 0, Score: 0}

	for _, seed := range winningSeeds {
		// result of 0 means incorrect pick / game not yet played
		if seed == 0 {
			continue
		}

		roundScore.Score += roundWeight * seed
		roundScore.Games += 1
	}

	return roundScore
}
