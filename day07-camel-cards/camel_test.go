package camel_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	camel "github.com/harveysanders/advent-of-code-2023/day07-camel-cards"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/harveysanders/advent-of-code-2023/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseHands(t *testing.T) {
	sample := strings.NewReader(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`)

	hands, err := camel.ParseHands(sample)
	require.NoError(t, err)

	require.Len(t, hands, 5)
	testCases := []struct {
		wantBid  int
		wantHand camel.Hand
	}{
		{
			wantHand: camel.Hand{
				Cards: [5]camel.Card{
					{Label: camel.Label3}, {Label: camel.Label2}, {Label: camel.LabelT}, {Label: camel.Label3}, {Label: camel.LabelK},
				},
			},
			wantBid: 765,
		},
		{
			wantHand: camel.Hand{
				Cards: [5]camel.Card{
					{Label: camel.LabelT}, {Label: camel.Label5}, {Label: camel.Label5}, {Label: camel.LabelJ}, {Label: camel.Label5},
				},
			},
			wantBid: 684,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("hand %d", i+1), func(t *testing.T) {
			hand := hands[i]
			require.Equal(t, tc.wantBid, hand.Bid)
			for j, gotCard := range hand.Cards {
				require.Equalf(t, tc.wantHand.Cards[j].Label, gotCard.Label, "card: %d", j+1)
			}
		})
	}
}

func TestHandType(t *testing.T) {
	testCases := []struct {
		labels   string
		wantType camel.HandType
	}{
		{labels: "32T3K", wantType: camel.OnePair},
		{labels: "KK677", wantType: camel.TwoPair},
		{labels: "KTJJT", wantType: camel.TwoPair},
		{labels: "T55J5", wantType: camel.ThreeOfAKind},
		{labels: "QQQJA", wantType: camel.ThreeOfAKind},
		{labels: "23456", wantType: camel.HighCard},
		{labels: "23332", wantType: camel.FullHouse},
		{labels: "AAAA9", wantType: camel.FourOfAKind},
	}

	for _, tc := range testCases {
		t.Run(tc.labels, func(t *testing.T) {
			hand := camel.Hand{}
			hand.ParseLabels(tc.labels)
			require.Equal(t, tc.wantType.String(), hand.Type(false).String())
		})
	}
}

func TestRank(t *testing.T) {
	sample := strings.NewReader(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
AQQQA 23
QKKKQ 13
`)
	testCases := []struct {
		input     io.Reader
		wantOrder []string // Ranked from lowest to highest
	}{
		{
			input: sample,
			wantOrder: []string{
				"32T3K", // one pair
				"KTJJT", // two pair
				"KK677", // two pair (K > T)
				"T55J5", // three of a kind
				"QQQJA", // three of a kind (Q > T)
				"QKKKQ",
				"AQQQA",
			},
		},
	}

	for _, tc := range testCases {
		label := strings.Join(tc.wantOrder, ",")
		t.Run(label, func(t *testing.T) {
			game := camel.NewGame()
			err := game.Parse(tc.input)
			require.NoError(t, err)
			game.Rank()

			got := []string{}
			for _, h := range game.Hands {
				got = append(got, h.String())
			}
			require.Equal(t, tc.wantOrder, got)
		})
	}
}

func TestWinnings(t *testing.T) {
	sample := strings.NewReader(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`)

	useLocal := os.Getenv("CI") == ""
	fullInput, err := github.GetInputFile(7, useLocal)
	require.NoError(t, err)

	defer fullInput.Close()

	testCases := []struct {
		name       string
		input      io.ReadSeeker
		wantTotal  int
		comparator testutil.Comparator // How to compare want, got. For example GREATER_THAN means got should be greater than the expected value.
		options    camel.GameOption
	}{
		{
			name:       "sample part 1",
			input:      sample,
			wantTotal:  6440,
			comparator: testutil.EQUAL,
		},
		{
			name:       "full input part 1",
			input:      fullInput,
			wantTotal:  248569531,
			comparator: testutil.EQUAL,
		},
		{
			name:       "sample part 2",
			input:      sample,
			wantTotal:  5905,
			comparator: testutil.EQUAL,
			options:    camel.WithWildcard(camel.LabelJ),
		},
		{
			name:       "full part 2",
			input:      fullInput,
			wantTotal:  253939737,
			comparator: testutil.EQUAL,
			options:    camel.WithWildcard(camel.LabelJ),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)
			game := camel.NewGame(tc.options)
			err = game.Parse(tc.input)
			require.NoError(t, err)

			gotTotal := game.TotalWinnings()
			switch tc.comparator {
			case testutil.EQUAL:
				require.Equal(t, tc.wantTotal, gotTotal)
			case testutil.GREATER_THAN:
				require.Greater(t, gotTotal, tc.wantTotal)
			case testutil.LESS_THAN:
				require.Less(t, gotTotal, tc.wantTotal)
			}
		})
	}
}
