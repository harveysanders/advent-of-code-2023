package cubes

import (
	"fmt"
	"strconv"
	"strings"
)

type Set struct {
	Red   int // Number of red cubes drawn.
	Green int // Number of green cubes drawn.
	Blue  int // Number of blue cubes drawn.
}

type Game struct {
	ID   int
	Sets []Set // A game has 3 sets
}

func (g *Game) Parse(data string) error {
	gameParts := strings.Split(data, ": ")
	if len(gameParts) != 2 {
		return fmt.Errorf("invalid game: %q", data)
	}

	// Parse ID
	rawID := strings.TrimPrefix(gameParts[0], "Game ")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w: %q", err, data)
	}

	g.ID = id

	// Parse Sets
	g.Sets = make([]Set, 0)

	rawSets := strings.Split(gameParts[1], "; ")
	for _, rawSet := range rawSets {
		set := Set{}
		rawCounts := strings.Split(rawSet, ", ")
		for _, rawCount := range rawCounts {
			p := strings.Fields(rawCount)
			if len(p) != 2 {
				return fmt.Errorf("invalid set: %q", rawSet)
			}
			rawCount := p[0]
			color := p[1]
			count, err := strconv.Atoi(rawCount)
			if err != nil {
				return fmt.Errorf("parse cube count: strconv.Atoi: %w", err)
			}
			switch color {
			case "red":
				set.Red = count
			case "green":
				set.Green = count
			case "blue":
				set.Blue = count
			}
		}
		g.Sets = append(g.Sets, set)
	}
	return nil
}

// A bag contains any positive number of red, green, and blue cubes.
type Bag struct {
	red   int // Number of red cubes.
	green int // Number of green cubes.
	blue  int // Number of blue cubes.
}

func NewBag(r, g, b int) *Bag {
	return &Bag{red: r, green: g, blue: b}
}

// ValidateGame returns true if the game is possible with the bag.
func (b *Bag) ValidateGame(game Game) bool {
	for _, s := range game.Sets {
		if s.Red > b.red {
			return false
		}
		if s.Green > b.green {
			return false
		}
		if s.Blue > b.blue {
			return false
		}
	}
	return true
}
