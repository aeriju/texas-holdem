package poker

import (
	"errors"
	"fmt"
	"strings"
)

type Card struct {
	Suit byte
	Rank byte
}

const (
	SuitClubs    = 'C'
	SuitDiamonds = 'D'
	SuitHearts   = 'H'
	SuitSpades   = 'S'
)

var rankToValue = map[byte]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var valueToRank = map[int]byte{
	2:  '2',
	3:  '3',
	4:  '4',
	5:  '5',
	6:  '6',
	7:  '7',
	8:  '8',
	9:  '9',
	10: 'T',
	11: 'J',
	12: 'Q',
	13: 'K',
	14: 'A',
}

func ParseCard(s string) (Card, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if len(s) != 2 {
		return Card{}, fmt.Errorf("invalid card '%s'", s)
	}
	suit := s[0]
	rank := s[1]
	if suit != SuitClubs && suit != SuitDiamonds && suit != SuitHearts && suit != SuitSpades {
		return Card{}, fmt.Errorf("invalid suit '%c'", suit)
	}
	if _, ok := rankToValue[rank]; !ok {
		return Card{}, fmt.Errorf("invalid rank '%c'", rank)
	}
	return Card{Suit: suit, Rank: rank}, nil
}

func ParseCards(list []string) ([]Card, error) {
	cards := make([]Card, 0, len(list))
	seen := map[string]struct{}{}
	for _, s := range list {
		c, err := ParseCard(s)
		if err != nil {
			return nil, err
		}
		key := c.String()
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicate card '%s'", key)
		}
		seen[key] = struct{}{}
		cards = append(cards, c)
	}
	return cards, nil
}

func (c Card) String() string {
	return fmt.Sprintf("%c%c", c.Suit, c.Rank)
}

func (c Card) RankValue() int {
	return rankToValue[c.Rank]
}

func RankToChar(v int) (byte, error) {
	r, ok := valueToRank[v]
	if !ok {
		return 0, errors.New("invalid rank value")
	}
	return r, nil
}

func NewDeck() []Card {
	suits := []byte{SuitClubs, SuitDiamonds, SuitHearts, SuitSpades}
	ranks := []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}
	deck := make([]Card, 0, 52)
	for _, s := range suits {
		for _, r := range ranks {
			deck = append(deck, Card{Suit: s, Rank: r})
		}
	}
	return deck
}

func RemoveCards(deck []Card, remove []Card) ([]Card, error) {
	toRemove := map[string]struct{}{}
	for _, c := range remove {
		toRemove[c.String()] = struct{}{}
	}
	filtered := make([]Card, 0, len(deck))
	for _, c := range deck {
		if _, ok := toRemove[c.String()]; !ok {
			filtered = append(filtered, c)
		}
	}
	if len(deck)-len(filtered) != len(toRemove) {
		return nil, errors.New("failed to remove cards from deck")
	}
	return filtered, nil
}
