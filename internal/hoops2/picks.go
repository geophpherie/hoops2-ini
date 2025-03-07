package hoops2

type Picks struct {
	Raw []int
}

type FinishingPicks struct {
	RegionFinalists []int
	FinalFour       []int // region winners
	FinalGame       []int // index to region winner
	Winner          int   // index to region winner
}

func NewPicks(picks []int) Picks {
	return Picks{Raw: picks}
}

func (p Picks) FinishingPicks() FinishingPicks {
	return FinishingPicks{
		RegionFinalists: p.RoundSeeds(3),
		FinalFour:       p.RoundSeeds(4),
		FinalGame:       p.FinalGame(),
		Winner:          p.Winner(),
	}
}
func (p Picks) FinalGame() []int { return p.Raw[60:62] }
func (p Picks) Winner() int      { return p.Raw[62] }
func (p Picks) RoundSeeds(roundNumber int) []int {
	// used in scoring
	switch roundNumber {
	case 1:
		return p.Raw[0:32]
	case 2:
		return p.Raw[32:48]
	case 3:
		return p.Raw[48:56]
	case 4:
		return p.Raw[56:60]
	case 5:
		var team1, team2 int
		if p.Raw[60] == 0 {
			team1 = 0
		} else {
			team1 = p.Raw[56:60][p.Raw[60]-1]
		}
		if p.Raw[61] == 0 {
			team2 = 0
		} else {
			team2 = p.Raw[56:60][p.Raw[61]-1]
		}
		return []int{
			team1,
			team2,
		}
	case 6:
		var team int
		if p.Raw[62] == 0 {
			team = 0
		} else {
			team = p.Raw[56:60][p.Raw[62]-1]
		}
		return []int{team}
	default:
		panic("invalid round number")
	}
}

func (p Picks) Region(regionNumber int) map[int][]int {
	switch regionNumber {
	case 1:
		return map[int][]int{
			1: p.Raw[0:8],
			2: p.Raw[32:36],
			3: p.Raw[48:50],
			4: p.Raw[56:57],
		}
	case 2:
		return map[int][]int{
			1: p.Raw[8:16],
			2: p.Raw[36:40],
			3: p.Raw[50:52],
			4: p.Raw[57:58],
		}
	case 3:
		return map[int][]int{
			1: p.Raw[16:24],
			2: p.Raw[40:44],
			3: p.Raw[52:54],
			4: p.Raw[58:59],
		}
	case 4:
		return map[int][]int{
			1: p.Raw[24:32],
			2: p.Raw[44:48],
			3: p.Raw[54:56],
			4: p.Raw[59:60],
		}
	default:
		panic("invalid region number")
	}
}

func (p Picks) Compare(other Picks) Picks {
	correctPicks := make([]int, len(p.Raw))
	for i := range len(p.Raw) {
		if p.Raw[i] == other.Raw[i] {
			correctPicks[i] = p.Raw[i]
		} else {
			correctPicks[i] = 0
		}
	}

	return Picks{Raw: correctPicks}
}
