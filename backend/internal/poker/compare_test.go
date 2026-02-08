package poker

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestComparison(t *testing.T) {
	cases, err := loadComparisonCases()
	if err != nil {
		t.Fatalf("load comparison cases: %v", err)
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			community := mustParseCards(t, tc.community)
			p1 := mustParseCards(t, tc.p1)
			p2 := mustParseCards(t, tc.p2)

			if len(community) != 5 || len(p1) != 2 || len(p2) != 2 {
				t.Fatalf("invalid card counts: community=%d p1=%d p2=%d", len(community), len(p1), len(p2))
			}

			h1 := append([]Card{}, p1...)
			h1 = append(h1, community...)
			h2 := append([]Card{}, p2...)
			h2 = append(h2, community...)

			h1Rank, err := Evaluate7(h1)
			if err != nil {
				t.Fatalf("evaluate hand1: %v", err)
			}
			h2Rank, err := Evaluate7(h2)
			if err != nil {
				t.Fatalf("evaluate hand2: %v", err)
			}

			got := Compare(h1Rank, h2Rank)
			want := expectedCmp(tc.result)
			if got != want {
				t.Fatalf("compare: got %d want %d", got, want)
			}
		})
	}
}

func mustParseCards(t *testing.T, s string) []Card {
	t.Helper()
	clean := strings.ReplaceAll(s, "\u00a0", " ")
	parts := strings.Fields(clean)
	cards, err := ParseCards(parts)
	if err != nil {
		t.Fatalf("parse cards %q: %v", s, err)
	}
	return cards
}

func expectedCmp(result int) int {
	switch result {
	case 1:
		return 1
	case 2:
		return -1
	default:
		return 0
	}
}

type compareCase struct {
	name      string
	community string
	p1        string
	p2        string
	result    int
}

func loadComparisonCases() ([]compareCase, error) {
	path, err := findComparisonCasesPath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty csv: %s", path)
	}

	var cases []compareCase
	for i, row := range records[1:] {
		if len(row) < 6 {
			continue
		}
		community := strings.TrimSpace(row[0])
		p1 := strings.TrimSpace(row[1])
		p2 := strings.TrimSpace(row[3])
		resultStr := strings.TrimSpace(row[5])

		if community == "" || p1 == "" || p2 == "" || resultStr == "" {
			continue
		}

		if len(strings.Fields(strings.ReplaceAll(community, "\u00a0", " "))) != 5 {
			continue
		}
		if len(strings.Fields(strings.ReplaceAll(p1, "\u00a0", " "))) != 2 {
			continue
		}
		if len(strings.Fields(strings.ReplaceAll(p2, "\u00a0", " "))) != 2 {
			continue
		}

		result, err := strconv.Atoi(resultStr)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid result %q", i+2, resultStr)
		}
		if result < 0 || result > 2 {
			return nil, fmt.Errorf("row %d: result must be 0,1,2", i+2)
		}

		cases = append(cases, compareCase{
			name:      fmt.Sprintf("row_%d", i+2),
			community: community,
			p1:        p1,
			p2:        p2,
			result:    result,
		})
	}

	if len(cases) == 0 {
		return nil, fmt.Errorf("no valid cases found in %s", path)
	}
	return cases, nil
}

func findComparisonCasesPath() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		try := filepath.Join(start, "assets", "test_cases", "comparison_test_cases.csv")
		if _, err := os.Stat(try); err == nil {
			return try, nil
		}
		start = filepath.Dir(start)
	}
	return "", fmt.Errorf("comparison_test_cases.csv not found")
}
