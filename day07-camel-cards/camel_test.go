package camel_test

import (
	"fmt"
	"strings"
	"testing"

	camel "github.com/harveysanders/advent-of-code-2023/day07-camel-cards"
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
