package scratchcards_test

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"

	scratchcards "github.com/harveysanders/advent-of-code-2023/day04-scratchcards"
	"github.com/stretchr/testify/require"
)

func TestParseCards(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`))

	cards, err := scratchcards.ParseCards(input)
	require.NoError(t, err)

	require.Len(t, cards, 6)
	card1 := cards[0]
	require.Equal(t, 1, card1.ID)
	wantWinning := []int{41, 48, 83, 86, 17}
	wantYours := []int{83, 86, 6, 31, 17, 9, 48, 53}
	require.Equal(t, wantWinning, card1.Winning)
	require.Equal(t, wantYours, card1.Yours)
}

func TestCardPoints(t *testing.T) {
	testCases := []struct {
		cardID     int
		wantPoints int
	}{
		{cardID: 1, wantPoints: 8},
		{cardID: 2, wantPoints: 2},
		{cardID: 3, wantPoints: 2},
		{cardID: 4, wantPoints: 1},
		{cardID: 5, wantPoints: 0},
		{cardID: 6, wantPoints: 0},
	}

	input := io.NopCloser(strings.NewReader(`Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`))

	cards, err := scratchcards.ParseCards(input)
	require.NoError(t, err)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Card %d", tc.cardID), func(t *testing.T) {
			card := cards[i]
			gotPoints := card.Points()
			require.Equal(t, tc.cardID, card.ID)
			require.Equal(t, tc.wantPoints, gotPoints)
		})
	}
}

//go:embed input/*
var inputFiles embed.FS

func TestCardsPoints(t *testing.T) {
	sampleInput := io.NopCloser(strings.NewReader(`Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`))

	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		name       string
		input      io.ReadCloser
		wantPoints int
	}{
		{
			name:       "Sample Part 1",
			input:      sampleInput,
			wantPoints: 13,
		},
		{
			name:       "Full input Part 1",
			input:      fullInput,
			wantPoints: 23941,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.input.Close()

			cards, err := scratchcards.ParseCards(tc.input)
			require.NoError(t, err)

			require.Equal(t, tc.wantPoints, cards.Points())
		})
	}
}

func TestCalcCopies(t *testing.T) {
	sampleInput := io.NopCloser(strings.NewReader(`Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`))

	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     io.ReadCloser
		wantTotal int
	}{
		{
			name:      "Part2 sample",
			input:     sampleInput,
			wantTotal: 30,
		},
		{
			name:      "Part2 full input",
			input:     fullInput,
			wantTotal: 5571760,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.input.Close()

			cards, err := scratchcards.ParseCards(tc.input)
			require.NoError(t, err)
			gotTotal := cards.CalcCopies()
			require.Equal(t, tc.wantTotal, gotTotal)
		})
	}
}
