package poker

import (
	"errors"
	"math/rand"
	"time"
)

func MonteCarlo(hole []Card, community []Card, players int, sims int) (float64, error) {
	if len(hole) != 2 {
		return 0, errors.New("hole must have 2 cards")
	}
	if !(len(community) == 0 || len(community) == 3 || len(community) == 4 || len(community) == 5) {
		return 0, errors.New("community must have 0, 3, 4, or 5 cards")
	}
	if players < 2 {
		return 0, errors.New("players must be >= 2")
	}
	if sims <= 0 {
		return 0, errors.New("simulations must be > 0")
	}

	known := append([]Card{}, hole...)
	known = append(known, community...)
	deck, err := RemoveCards(NewDeck(), known)
	if err != nil {
		return 0, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	wins := 0.0

	for i := 0; i < sims; i++ {
		simDeck := append([]Card{}, deck...)
		rng.Shuffle(len(simDeck), func(i, j int) { simDeck[i], simDeck[j] = simDeck[j], simDeck[i] })
		idx := 0

		neededCommunity := 5 - len(community)
		board := append([]Card{}, community...)
		board = append(board, simDeck[idx:idx+neededCommunity]...)
		idx += neededCommunity

		hero7 := append([]Card{}, hole...)
		hero7 = append(hero7, board...)
		heroRank, err := Evaluate7(hero7)
		if err != nil {
			return 0, err
		}

		bestOpp := heroRank
		winners := 1

		for p := 0; p < players-1; p++ {
			oppHole := []Card{simDeck[idx], simDeck[idx+1]}
			idx += 2
			opp7 := append([]Card{}, oppHole...)
			opp7 = append(opp7, board...)
			oppRank, err := Evaluate7(opp7)
			if err != nil {
				return 0, err
			}
			cmp := Compare(oppRank, bestOpp)
			if cmp > 0 {
				bestOpp = oppRank
				winners = 1
			} else if cmp == 0 {
				winners++
			}
		}

		if Compare(heroRank, bestOpp) == 0 {
			wins += 1.0 / float64(winners)
		}
	}

	return wins / float64(sims), nil
}
