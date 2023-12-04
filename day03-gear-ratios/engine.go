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

type Number struct {
	Value    int
	Location Coord
	Size     int
}

type Schematic struct {
	matrix []string
	maxX   int
	maxY   int
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

	s.maxY = len(s.matrix)
	if s.maxY > 0 {
		s.maxX = len(s.matrix[0])
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

func (s *Schematic) IsPartNum(n Number) bool {
	leftAdjX := int(math.Max(float64(n.Location.X-1), 0))
	rightAdjX := int(math.Min(float64(n.Location.X+n.Size), float64(s.maxX)))
	topAdjY := int(math.Max(float64(n.Location.Y-1), 0))
	bottomAdjY := int(math.Min(float64(n.Location.Y+1), float64(s.maxY)))

	symRE := regexp.MustCompile(`[^\d|\.]`)
	if n.Location.Y != 0 {
		upperLine := s.matrix[topAdjY][leftAdjX : rightAdjX+1]
		if symRE.MatchString(upperLine) {
			return true
		}
	}

	if n.Location.Y != s.maxY {
		bottomLine := s.matrix[bottomAdjY][leftAdjX : rightAdjX+1]
		if symRE.MatchString(bottomLine) {
			return true
		}
	}

	if n.Location.X != 0 {
		char := s.matrix[n.Location.Y][leftAdjX : leftAdjX+1]
		if symRE.MatchString(char) {
			return true
		}
	}

	if n.Location.X+n.Size != s.maxX {
		char := s.matrix[n.Location.Y][rightAdjX : rightAdjX+1]
		if symRE.MatchString(char) {
			return true
		}
	}
	return false
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
