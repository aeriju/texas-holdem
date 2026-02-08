package poker

import "testing"

func TestParseCardValidation(t *testing.T) {
	if _, err := ParseCard("HA"); err != nil {
		t.Fatalf("expected valid card, got %v", err)
	}
	if _, err := ParseCard("X9"); err == nil {
		t.Fatalf("expected invalid suit error")
	}
	if _, err := ParseCard("H1"); err == nil {
		t.Fatalf("expected invalid rank error")
	}
	if _, err := ParseCard("H10"); err == nil {
		t.Fatalf("expected invalid length error")
	}
	_, err := ParseCards([]string{"HA", "HA"})
	if err == nil {
		t.Fatalf("expected duplicate card error")
	}
}

func TestEvaluate7Categories(t *testing.T) {
	tests := []struct {
		name     string
		cards    []string
		category Category
	}{
		{
			name:     "straight flush",
			cards:    []string{"H9", "HT", "HJ", "HQ", "HK", "C2", "D3"},
			category: StraightFlush,
		},
		{
			name:     "four of a kind",
			cards:    []string{"S9", "H9", "D9", "C9", "HK", "C2", "D3"},
			category: FourOfAKind,
		},
		{
			name:     "full house",
			cards:    []string{"S9", "H9", "D9", "CK", "HK", "C2", "D3"},
			category: FullHouse,
		},
		{
			name:     "flush",
			cards:    []string{"H2", "H5", "H7", "HJ", "HK", "C2", "D3"},
			category: Flush,
		},
		{
			name:     "straight",
			cards:    []string{"H9", "CT", "HJ", "DQ", "CK", "C2", "D3"},
			category: Straight,
		},
		{
			name:     "three of a kind",
			cards:    []string{"S9", "H9", "D9", "CK", "HQ", "C2", "D3"},
			category: ThreeOfAKind,
		},
		{
			name:     "two pair",
			cards:    []string{"S9", "H9", "DK", "CK", "HQ", "C2", "D3"},
			category: TwoPair,
		},
		{
			name:     "one pair",
			cards:    []string{"S9", "H9", "DK", "CQ", "HJ", "C2", "D3"},
			category: OnePair,
		},
		{
			name:     "high card",
			cards:    []string{"S9", "H7", "DK", "CQ", "HJ", "C2", "D3"},
			category: HighCard,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ParseCards(tc.cards)
			if err != nil {
				t.Fatalf("parse cards: %v", err)
			}
			rank, err := Evaluate7(parsed)
			if err != nil {
				t.Fatalf("evaluate: %v", err)
			}
			if rank.Category != tc.category {
				t.Fatalf("category: got %v want %v", rank.Category, tc.category)
			}
		})
	}
}

func TestMonteCarloValidation(t *testing.T) {
	if _, err := MonteCarlo([]Card{}, []Card{}, 2, 100); err == nil {
		t.Fatalf("expected hole size error")
	}
	if _, err := MonteCarlo([]Card{{Suit: 'H', Rank: 'A'}, {Suit: 'S', Rank: 'K'}}, []Card{}, 1, 100); err == nil {
		t.Fatalf("expected players error")
	}
	if _, err := MonteCarlo([]Card{{Suit: 'H', Rank: 'A'}, {Suit: 'S', Rank: 'K'}}, []Card{}, 2, 0); err == nil {
		t.Fatalf("expected simulations error")
	}
	if _, err := MonteCarlo([]Card{{Suit: 'H', Rank: 'A'}, {Suit: 'S', Rank: 'K'}}, []Card{{Suit: 'C', Rank: '2'}}, 2, 100); err == nil {
		t.Fatalf("expected community size error")
	}
}
