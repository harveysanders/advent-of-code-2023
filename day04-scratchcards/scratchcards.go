package scratchcards

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"
)

type ScratchCard struct {
	ID      int   // The ID of the card.
	Winning []int // List of winning numbers, i.e. if your any of your numbers match a winning number, you win points.
	Yours   []int // Your pre-selected numbers. If any of these numbers match a winning number, you win!
}

func (s *ScratchCard) Decode(raw string) error {
	p := strings.Split(raw, ": ")
	if len(p) != 2 {
		return fmt.Errorf("invalid card data: %q", raw)
	}
	rawID := strings.TrimSpace(strings.TrimPrefix(p[0], "Card"))
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return fmt.Errorf("decode ID: %w, val: %q", err, rawID)
	}
	s.ID = id

	rawNums := strings.Split(p[1], " | ")
	if len(rawNums) != 2 {
		return fmt.Errorf("invalid card numbers: %q", p[1])
	}

	s.Winning = make([]int, 0)
	s.Yours = make([]int, 0)
	for i, r := range rawNums {
		nums := strings.Fields(r)
		for _, rawN := range nums {
			n, err := strconv.Atoi(strings.TrimSpace(rawN))
			if err != nil {
				return fmt.Errorf("invalid number: %w, val: %q", err, rawN)
			}
			switch i {
			case 0:
				s.Winning = append(s.Winning, n)
			case 1:
				s.Yours = append(s.Yours, n)
			}
		}
	}
	return nil
}

func (s ScratchCard) Points() int {
	winners := make([]int, 0)
	for _, yourN := range s.Yours {
		if !slices.Contains(s.Winning, yourN) {
			continue
		}
		winners = append(winners, yourN)
	}
	return int(math.Pow(2, float64(len(winners)-1)))
}

type Cards []ScratchCard

func ParseCards(r io.Reader) (Cards, error) {
	cards := make(Cards, 0)

	scr := bufio.NewScanner(r)
	for scr.Scan() {
		if scr.Err() != nil {
			return cards, fmt.Errorf("scr.Error(): %w", scr.Err())
		}

		line := scr.Text()
		card := ScratchCard{}
		err := card.Decode(line)
		if err != nil {
			return cards, fmt.Errorf("card.Decode: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (cs Cards) Points() int {
	sum := 0
	for _, c := range cs {
		sum += c.Points()
	}
	return sum
}
