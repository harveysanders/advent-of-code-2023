package camel

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math"
	"slices"
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
	HighCard     HandType = iota // HighCard is where all cards' labels are distinct: 23456
	OnePair                      // OnePair is where two cards share one label, and the other three cards have a different label from the pair and each other: A23A4
	TwoPair                      // TwoPair is where two cards share one label, two other cards share a second label, and the remaining card has a third label: 23432
	ThreeOfAKind                 // ThreeOfAKind is where three cards have the same label, and the remaining two cards are each different from any other card in the hand: TTT98
	FullHouse                    // A full house is where three cards have the same label, and the remaining two cards share a different label: 23332
	FourOfAKind                  // FourOfAKind is where four cards have the same label and one card has a different label: AA8AA
	FiveOfAKind                  // FiveOfAKind is where all five cards have the same label: AAAAA.
)

func (t HandType) String() string {
	names := []string{"", "high card", "one pair", "two pair", "three of a kind", "four of a kind", "five of a kind"}
	return names[t]
}

type Game struct {
	wildcard Label
	Hands    Hands
}

type GameOption func(*Game)

func WithWildcard(card Label) GameOption {
	return func(g *Game) {
		g.wildcard = card
	}
}

func NewGame(opts ...GameOption) *Game {
	g := &Game{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(g)
	}
	return g
}

func (g *Game) Parse(r io.Reader) error {
	hands, err := ParseHands(r, withGame(g))
	if err != nil {
		return err
	}
	g.Hands = hands
	return nil
}

type Card struct {
	Label Label
	value int
}

type Hand struct {
	Cards [5]Card
	Bid   int
	game  *Game
}

// Type finds the number of sets of matching cards and returns the associated HandType.
func (h Hand) Type(useWildcard bool) HandType {
	counts := h.cardCounts()
	pairs := []Label{}
	maxMatches := 0
	wildCardCount := 0
	if useWildcard && h.game != nil {
		wildCardCount = counts[h.game.wildcard]
	}
	for label, count := range counts {
		maxMatches = int(math.Max(float64(maxMatches+wildCardCount), float64(count)))
		if count+wildCardCount == 2 {
			pairs = append(pairs, label)
		}
	}
	// If 4 of a kind or higher, use that for type
	if maxMatches >= 4 {
		return HandType(maxMatches + 1)
	}

	// # of pairs in the hand
	nPairs := len(pairs)
	if maxMatches == 3 {
		if nPairs == 1 {
			return FullHouse
		}
		return ThreeOfAKind
	}
	if nPairs > 0 {
		return HandType(nPairs)
	}
	// No pairs, all distinct
	return HighCard
}

// CardCounts returns a map of card labels to their counts in the hand.
// Ex:
//
//	"JJQQ3" -> {"J":2, "Q": 2, "3": 1}
func (h Hand) cardCounts() map[Label]int {
	cardCounts := map[Label]int{}
	for _, card := range h.Cards {
		count, ok := cardCounts[card.Label]
		if !ok {
			cardCounts[card.Label] = 1
		}
		cardCounts[card.Label] = count + 1
	}
	return cardCounts
}

func (h Hand) String() string {
	var s strings.Builder
	for _, c := range h.Cards {
		if _, err := s.WriteString(string(c.Label)); err != nil {
			fmt.Printf("s.WriteString(): %v, val: %q", err, c.Label)
			return ""
		}
	}
	return s.String()
}

func (h Hand) hasWildcard() int {
	if h.game == nil {
		return 0
	}
	if h.game.wildcard == "" {
		return 0
	}
	if strings.Contains(h.String(), string(h.game.wildcard)) {
		return 1
	}
	return 0
}

// ParseLabels takes a hand as a string, e.g. "K234J" and populates the hand struct.
func (h *Hand) ParseLabels(raw string) {
	labels := strings.Split(raw, "")
	for i, l := range labels {
		label := Label(l)
		value := cardValues[label]
		if h.game != nil && h.game.wildcard == label {
			value = 1
		}
		h.Cards[i] = Card{
			Label: Label(label),
			value: value,
		}
	}
}

type Hands []Hand

// Rank sorts the hands in place, ordering by type strength, lowest rank first.
func (g Game) Rank() Hands {
	copy := slices.Clone(g.Hands)
	slices.SortStableFunc(copy, cmpHands(g.wildcard != ""))
	return copy
}

func cmpHands(useWildcard bool) func(a Hand, b Hand) int {
	return func(a, b Hand) int {
		if n := cmp.Compare(a.Type(useWildcard), b.Type(useWildcard)); n != 0 {
			return n
		}

		if useWildcard {
			if n := cmpHands(false)(a, b); n != 0 {
				return n
			}
		}

		// If types are equal, order by card value
		for i, aCard := range a.Cards {
			if n := cmp.Compare(aCard.value, b.Cards[i].value); n != 0 {
				return n
			}
		}
		return 0
	}
}

// TotalWinnings ranks the hands by type, then calculates the winnings based on the hand's bid and rank.
func (g Game) TotalWinnings() int {
	ranked := g.Rank()
	var total float64
	for i, hand := range ranked {
		rank := i + 1
		handWinnings := hand.Bid * rank
		total += float64(handWinnings)
	}
	return int(total)
}

type handOption func(h *Hand)

func withGame(g *Game) handOption {
	return func(h *Hand) {
		h.game = g
	}
}

func ParseHands(r io.Reader, opts ...handOption) (Hands, error) {
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

		rawBid := p[1]
		bid, err := strconv.Atoi(rawBid)
		if err != nil {
			return hands, fmt.Errorf("strconv.Atoi: %w, val: %q", err, rawBid)
		}

		hand := Hand{
			Cards: [5]Card{},
			Bid:   bid,
		}
		for _, o := range opts {
			o(&hand)
		}

		hand.ParseLabels(p[0])
		hands = append(hands, hand)
	}

	return hands, nil
}
