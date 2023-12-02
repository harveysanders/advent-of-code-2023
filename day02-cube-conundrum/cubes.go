package cubes

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

// FewestCubes returns the fewest possible amount of cubes for each color to play the given game.
func (g *Game) FewestCubes() (red, green, blue int) {
	for _, s := range g.Sets {
		red = int(math.Max(float64(s.Red), float64(red)))
		green = int(math.Max(float64(s.Green), float64(green)))
		blue = int(math.Max(float64(s.Blue), float64(blue)))
	}
	return red, green, blue
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

type Record struct {
	games []Game
}

func (r *Record) Decode(rdr io.Reader) error {
	r.games = make([]Game, 0)

	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		if scr.Err() != nil {
			return fmt.Errorf("scr.Scan(): %w", scr.Err())
		}

		line := scr.Text()
		game := Game{}
		if err := game.Parse(line); err != nil {
			return fmt.Errorf("game.Parse(): %w - line: %q", err, line)
		}
		r.games = append(r.games, game)
	}
	return nil
}

// ValidGameIDs returns a list of valid game IDs for a given bag.
func (r *Record) ValidGameIDs(b Bag) []int {
	results := make([]int, 0)
	for _, g := range r.games {
		if b.ValidateGame(g) {
			results = append(results, g.ID)
		}
	}
	return results
}

// Part2 calculated the fewest possible cubes for each game, calculates the power of those counts for each game, then returns the sum of all the powers.
func (r *Record) Part2() int {
	sum := 0
	for _, game := range r.games {
		red, green, blue := game.FewestCubes()
		power := red * green * blue
		sum += power
	}
	return sum
}
