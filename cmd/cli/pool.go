package main

type Pool struct {
	Year         string
	RoundWeights map[int]int
	Regions      map[int]Region
	Results      []int
	Contestants  map[string]Contestant
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

	pool.Year = filepath.Ext(INI_FILE)[2:]

	return pool, nil
}

func (p Pool) ScoreTruth() int {
	round1 := p.ScoreRound(p.Results[0:32], 1, nil)
	round1.Display()

	round2 := p.ScoreRound(p.Results[32:48], 2, &round1)
	round2.Display()

	round3 := p.ScoreRound(p.Results[48:56], 3, &round2)
	round3.Display()

	round4 := p.ScoreRound(p.Results[56:60], 4, &round3)
	round4.Display()

	// need different way to get teams here
	// p.results[60] and p.results[61] are index into p.Results[56:60]
	round5 := p.ScoreRound([]int{
		p.Results[56:60][p.Results[60]-1],
		p.Results[56:60][p.Results[61]-1],
	}, 5, &round4)
	round5.Display()

	// p.results[66] is index into p.Results[56:60]
	round6 := p.ScoreRound([]int{
		p.Results[56:60][p.Results[62]-1],
	}, 6, &round5)
	round6.Display()

	return 0
}

func (p Pool) ScoreRound(results []int, roundNumber int, prevRound *RoundScore) RoundScore {
	score := 0
	games := 0
	for _, result := range results {
		score += p.RoundWeights[roundNumber] * result
		if result != 0 {
			games += 1
		}
	}

	roundScore := RoundScore{Number: roundNumber, Games: games, Score: score}
	if prevRound != nil {
		roundScore.Add(prevRound)
	}

	return roundScore
}

type Region struct {
	Name   string
	Number int
	Teams  map[int]string
}

type Contestant struct {
	Name  string
	Picks []int
}
