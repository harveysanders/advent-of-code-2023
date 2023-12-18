package image

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Bit string

const (
	emptySpace Bit = "."
	galaxyBit  Bit = "#"
)

// Observation represents the pixel data of an observatory image. "." represents empty space, and "#" represents a galaxy.
type Observation struct {
	rows []string
}

func (o Observation) Width() int {
	if len(o.rows) == 0 {
		return 0
	}
	return len(o.rows[0])
}

func (o Observation) Height() int {
	return len(o.rows)
}

func (o Observation) columnAt(x int) string {
	var s strings.Builder
	for y := 0; y < o.Height(); y++ {
		s.WriteByte(o.rows[y][x])
	}
	return s.String()
}

func (o Observation) String() string {
	return strings.Join(o.rows, "\n")
}

func ParseImage(r io.Reader) (Observation, error) {
	scr := bufio.NewScanner(r)
	d := Observation{rows: []string{}}
	for scr.Scan() {
		if scr.Err() != nil {
			return d, fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		d.rows = append(d.rows, scr.Text())
	}

	return d, nil
}

// Expand creates a new observation image where any rows or columns in the original image that contain no galaxies are doubled in size.
func (o Observation) Expand() Observation {
	expanded := Observation{rows: []string{}}

	for _, row := range o.rows {
		if strings.Contains(row, string(galaxyBit)) {
			expanded.rows = append(expanded.rows, row)
			continue
		}
		// All empty space, double the row
		expanded.rows = append(expanded.rows, row, row)
	}

	// Track the # of times a column is doubled so we can calculate the correct X value
	doubledColN := 0
	for x := 0; x < o.Width(); x++ {
		col := o.columnAt(x)
		if strings.Contains(col, string(galaxyBit)) {
			continue
		}
		// All empty space, double the column
		for y, row := range expanded.rows {
			expanded.rows[y] = string(slices.Insert([]byte(row), x+doubledColN, '.'))
		}
		doubledColN++
	}

	return expanded
}
