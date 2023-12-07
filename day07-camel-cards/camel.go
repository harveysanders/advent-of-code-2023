package camel

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Label string

const (
	LabelA Label = "A"
	LabelK Label = "K"
	LabelQ Label = "Q"
	LabelJ Label = "J"
	LabelT Label = "T"
	Label9 Label = "9"
	Label8 Label = "8"
	Label7 Label = "7"
	Label6 Label = "6"
	Label5 Label = "5"
	Label4 Label = "4"
	Label3 Label = "3"
	Label2 Label = "2"
)

var cardValues = map[Label]int{
	LabelA: 14,
	LabelK: 13,
	LabelQ: 12,
	LabelJ: 11,
	LabelT: 10,
	Label9: 9,
	Label8: 8,
	Label7: 7,
	Label6: 6,
	Label5: 5,
	Label4: 4,
	Label3: 3,
	Label2: 2,
}

type HandType int

const (
	FiveOfAKind  HandType = iota // FiveOfAKind is where all five cards have the same label: AAAAA.
	FourOfAKind                  // FourOfAKind is where four cards have the same label and one card has a different label: AA8AA
	ThreeOfAKind                 // ThreeOfAKind is where three cards have the same label, and the remaining two cards are each different from any other card in the hand: TTT98
	TwoPair                      // TwoPair is where two cards share one label, two other cards share a second label, and the remaining card has a third label: 23432
	OnePair                      // OnePair is where two cards share one label, and the other three cards have a different label from the pair and each other: A23A4
	HighCard                     // HighCard is where all cards' labels are distinct: 23456
)

type Card struct {
	Label Label
	value int
}

type Hand struct {
	Cards [5]Card
	Type  HandType
	Bid   int
}

func ParseHands(r io.Reader) ([]Hand, error) {
	hands := make([]Hand, 0)
	scr := bufio.NewScanner(r)

	for scr.Scan() {
		if scr.Err() != nil {
			return hands, fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		line := scr.Text()
		p := strings.Fields(line)
		if len(p) != 2 {
			return hands, fmt.Errorf("invalid hand-bid line: %q", line)
		}

		labels := strings.Split(p[0], "")
		rawBid := p[1]
		bid, err := strconv.Atoi(rawBid)
		if err != nil {
			return hands, fmt.Errorf("strconv.Atoi: %w, val: %q", err, rawBid)
		}

		hand := Hand{
			Cards: [5]Card{},
			Bid:   bid,
		}
		for i, l := range labels {
			label := Label(l)
			hand.Cards[i] = Card{
				Label: Label(label),
				value: cardValues[label],
			}
		}
		hands = append(hands, hand)
	}

	return hands, nil
}
