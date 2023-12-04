package engine

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
)

type Coord struct {
	X int
	Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

type Number struct {
	Value    int
	Location Coord
	Size     int
}

type Schematic struct {
	matrix []string
	width  int
	height int
	gears  map[string][]Number
}

func (s *Schematic) Parse(r io.ReadCloser) error {
	s.matrix = make([]string, 0)
	scr := bufio.NewScanner(r)
	for scr.Scan() {
		if scr.Err() != nil {
			return fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		line := scr.Text()
		s.matrix = append(s.matrix, line)
	}

	s.height = len(s.matrix)
	if s.height > 0 {
		s.width = len(s.matrix[0])
	}
	return nil
}

func (s *Schematic) CollectNumbers() ([]Number, error) {
	numRe := regexp.MustCompile(`(\d+)`)
	res := make([]Number, 0)
	for y, row := range s.matrix {
		matches := numRe.FindAllStringIndex(row, -1)
		for _, matchLoc := range matches {
			x := matchLoc[0]
			asciiN := row[matchLoc[0]:matchLoc[1]]
			val, err := strconv.Atoi(string(asciiN))
			if err != nil {
				return res, fmt.Errorf("strconv.Atoi: %w", err)
			}
			n := Number{
				Value:    val,
				Location: Coord{X: x, Y: y},
				Size:     matchLoc[1] - matchLoc[0],
			}
			res = append(res, n)
		}
	}
	return res, nil
}

func (s *Schematic) hasAdjacentSymbol(n Number, symRE regexp.Regexp) (bool, Coord) {
	leftAdjX := int(math.Max(float64(n.Location.X-1), 0))
	rightAdjX := int(math.Min(float64(n.Location.X+n.Size), float64(s.width)))
	topAdjY := int(math.Max(float64(n.Location.Y-1), 0))
	bottomAdjY := int(math.Min(float64(n.Location.Y+1), float64(s.height)))

	rightEdge := math.Min(float64(rightAdjX+1), float64(s.width-1))
	if n.Location.Y != 0 {
		upperLine := s.matrix[topAdjY][leftAdjX:int(rightEdge)]
		loc := symRE.FindStringIndex(upperLine)
		if loc != nil {
			return true, Coord{X: loc[0] + leftAdjX, Y: topAdjY}
		}
	}

	if n.Location.Y != s.height-1 {
		bottomLine := s.matrix[bottomAdjY][leftAdjX:int(rightEdge)]
		loc := symRE.FindStringIndex(bottomLine)
		if loc != nil {
			return true, Coord{X: loc[0] + leftAdjX, Y: bottomAdjY}
		}
	}

	if n.Location.X != 0 {
		char := s.matrix[n.Location.Y][leftAdjX : leftAdjX+1]
		if symRE.MatchString(char) {
			return true, Coord{X: leftAdjX, Y: n.Location.Y}
		}
	}

	if n.Location.X+n.Size != s.width {
		char := s.matrix[n.Location.Y][rightAdjX:int(rightEdge)]
		if symRE.MatchString(char) {
			return true, Coord{X: rightAdjX, Y: n.Location.Y}
		}
	}
	return false, Coord{}
}

func (s *Schematic) IsPartNum(n Number) bool {
	symRE := regexp.MustCompile(`[^\d\.]`)
	isP, _ := s.hasAdjacentSymbol(n, *symRE)
	return isP
}

func (s *Schematic) PartNums() ([]int, error) {
	partNums := make([]int, 0)
	allNums, err := s.CollectNumbers()
	if err != nil {
		return partNums, fmt.Errorf("s.CollectNumbers(): %w", err)
	}

	for _, n := range allNums {
		if s.IsPartNum(n) {
			partNums = append(partNums, n.Value)
		}
	}
	return partNums, nil
}

func (s *Schematic) PartNumSum() (int, error) {
	nums, err := s.PartNums()
	if err != nil {
		return 0, fmt.Errorf("s.PartNums(): %w", err)
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum, nil
}

func (s *Schematic) FindGears() (int, error) {
	re := regexp.MustCompile(`\*`)

	allNums, err := s.CollectNumbers()
	if err != nil {
		return 0, err
	}

	s.gears = make(map[string][]Number)
	for _, n := range allNums {
		isGear, loc := s.hasAdjacentSymbol(n, *re)
		if isGear {
			key := loc.String()
			g, ok := s.gears[key]
			if !ok {
				s.gears[key] = []Number{n}
				continue
			}
			s.gears[key] = append(g, n)
		}
	}

	sum := 0
	for _, g := range s.gears {
		if len(g) == 2 {
			sum += g[0].Value * g[1].Value
		}
	}
	return sum, nil
}
