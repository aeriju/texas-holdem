package poker

import (
	"errors"
	"sort"
)

type Category int

const (
	HighCard Category = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

var categoryName = map[Category]string{
	HighCard:      "high card",
	OnePair:       "one pair",
	TwoPair:       "two pair",
	ThreeOfAKind:  "three of a kind",
	Straight:      "straight",
	Flush:         "flush",
	FullHouse:     "full house",
	FourOfAKind:   "four of a kind",
	StraightFlush: "straight flush",
}

type HandRank struct {
	Category Category
	Tiebreak []int
	Best5    []Card
}

func (h HandRank) Name() string {
	return categoryName[h.Category]
}

// computes the best 5-card hand from exactly 7 cards.
func Evaluate7(cs []Card) (HandRank, error) {
	if len(cs) != 7 {
		return HandRank{}, errors.New("Evaluate7 expects 7 cards")
	}
	best := HandRank{Category: -1}
	combos := combinations7to5(cs)
	for _, combo := range combos {
		r := Evaluate5(combo)
		if best.Category == -1 || Compare(r, best) > 0 {
			best = r
		}
	}
	return best, nil
}

// ranks a 5-card hand and returns its category and tiebreakers.
func Evaluate5(cs []Card) HandRank {
	ranks := make([]int, 0, 5)
	suits := make(map[byte]int)
	counts := make(map[int]int)
	for _, c := range cs {
		v := c.RankValue()
		ranks = append(ranks, v)
		counts[v]++
		suits[c.Suit]++
	}
	sort.Ints(ranks)
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] > ranks[j] })

	isFlush := false
	for _, n := range suits {
		if n == 5 {
			isFlush = true
			break
		}
	}

	isStraight, straightHigh := straightHigh(ranks)

	if isStraight && isFlush {
		return HandRank{
			Category: StraightFlush,
			Tiebreak: []int{straightHigh},
			Best5:    append([]Card{}, cs...),
		}
	}

	if rank, kicker, ok := fourOfAKind(counts); ok {
		return HandRank{
			Category: FourOfAKind,
			Tiebreak: []int{rank, kicker},
			Best5:    append([]Card{}, cs...),
		}
	}

	if trip, pair, ok := fullHouse(counts); ok {
		return HandRank{
			Category: FullHouse,
			Tiebreak: []int{trip, pair},
			Best5:    append([]Card{}, cs...),
		}
	}

	if isFlush {
		return HandRank{
			Category: Flush,
			Tiebreak: append([]int{}, ranks...),
			Best5:    append([]Card{}, cs...),
		}
	}

	if isStraight {
		return HandRank{
			Category: Straight,
			Tiebreak: []int{straightHigh},
			Best5:    append([]Card{}, cs...),
		}
	}

	if trip, kickers, ok := threeOfAKind(counts); ok {
		return HandRank{
			Category: ThreeOfAKind,
			Tiebreak: append([]int{trip}, kickers...),
			Best5:    append([]Card{}, cs...),
		}
	}

	if highPair, lowPair, kicker, ok := twoPair(counts); ok {
		return HandRank{
			Category: TwoPair,
			Tiebreak: []int{highPair, lowPair, kicker},
			Best5:    append([]Card{}, cs...),
		}
	}

	if pair, kickers, ok := onePair(counts); ok {
		return HandRank{
			Category: OnePair,
			Tiebreak: append([]int{pair}, kickers...),
			Best5:    append([]Card{}, cs...),
		}
	}

	return HandRank{
		Category: HighCard,
		Tiebreak: append([]int{}, ranks...),
		Best5:    append([]Card{}, cs...),
	}
}

// compares two ranked hands: 1 if a > b, -1 if a < b, 0 if equal.
func Compare(a, b HandRank) int {
	if a.Category != b.Category {
		if a.Category > b.Category {
			return 1
		}
		return -1
	}
	limit := len(a.Tiebreak)
	if len(b.Tiebreak) < limit {
		limit = len(b.Tiebreak)
	}
	for i := 0; i < limit; i++ {
		if a.Tiebreak[i] > b.Tiebreak[i] {
			return 1
		}
		if a.Tiebreak[i] < b.Tiebreak[i] {
			return -1
		}
	}
	return 0
}

// determines if the ranks form a straight and returns its high card.
func straightHigh(ranks []int) (bool, int) {
	unique := make([]int, 0, 5)
	seen := map[int]struct{}{}
	for _, r := range ranks {
		if _, ok := seen[r]; !ok {
			unique = append(unique, r)
			seen[r] = struct{}{}
		}
	}
	sort.Slice(unique, func(i, j int) bool { return unique[i] > unique[j] })
	if len(unique) != 5 {
		return false, 0
	}
	if unique[0]-unique[4] == 4 {
		return true, unique[0]
	}
	if unique[0] == 14 && unique[1] == 5 && unique[2] == 4 && unique[3] == 3 && unique[4] == 2 {
		return true, 5
	}
	return false, 0
}

// returns the quad rank and kicker rank if present.
func fourOfAKind(counts map[int]int) (int, int, bool) {
	quad := 0
	kicker := 0
	for r, c := range counts {
		if c == 4 {
			quad = r
		} else if c == 1 {
			kicker = r
		}
	}
	if quad == 0 {
		return 0, 0, false
	}
	return quad, kicker, true
}

// returns trip and pair ranks if present.
func fullHouse(counts map[int]int) (int, int, bool) {
	trip := 0
	pair := 0
	for r, c := range counts {
		if c == 3 && r > trip {
			trip = r
		}
	}
	if trip == 0 {
		return 0, 0, false
	}
	for r, c := range counts {
		if c >= 2 && r != trip {
			if r > pair {
				pair = r
			}
		}
	}
	if pair == 0 {
		return 0, 0, false
	}
	return trip, pair, true
}

// returns trip rank and kickers if present.
func threeOfAKind(counts map[int]int) (int, []int, bool) {
	trip := 0
	for r, c := range counts {
		if c == 3 && r > trip {
			trip = r
		}
	}
	if trip == 0 {
		return 0, nil, false
	}
	kickers := make([]int, 0, 2)
	for r, c := range counts {
		if c == 1 {
			kickers = append(kickers, r)
		}
	}
	sort.Slice(kickers, func(i, j int) bool { return kickers[i] > kickers[j] })
	return trip, kickers, true
}

// returns high pair, low pair, and kicker if present.
func twoPair(counts map[int]int) (int, int, int, bool) {
	pairs := make([]int, 0, 2)
	kicker := 0
	for r, c := range counts {
		if c == 2 {
			pairs = append(pairs, r)
		} else if c == 1 {
			kicker = r
		}
	}
	if len(pairs) != 2 {
		return 0, 0, 0, false
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i] > pairs[j] })
	return pairs[0], pairs[1], kicker, true
}

// returns pair rank and kickers if present.
func onePair(counts map[int]int) (int, []int, bool) {
	pair := 0
	for r, c := range counts {
		if c == 2 && r > pair {
			pair = r
		}
	}
	if pair == 0 {
		return 0, nil, false
	}
	kickers := make([]int, 0, 3)
	for r, c := range counts {
		if c == 1 {
			kickers = append(kickers, r)
		}
	}
	sort.Slice(kickers, func(i, j int) bool { return kickers[i] > kickers[j] })
	return pair, kickers, true
}

// enumerates all 21 combinations of 5 cards from 7.
func combinations7to5(cs []Card) [][]Card {
	combos := make([][]Card, 0, 21)
	for i := 0; i < 7; i++ {
		for j := i + 1; j < 7; j++ {
			combo := make([]Card, 0, 5)
			for k := 0; k < 7; k++ {
				if k == i || k == j {
					continue
				}
				combo = append(combo, cs[k])
			}
			combos = append(combos, combo)
		}
	}
	return combos
}
