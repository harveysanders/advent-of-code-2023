package mirror

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

func (o Orientation) String() string {
	return []string{"Horizontal", "Vertical"}[o]
}

type Pattern struct {
	grid []string
}

func (p Pattern) height() int {
	return len(p.grid)
}
func (p Pattern) width() int {
	if len(p.grid) == 0 {
		return 0
	}
	return len(p.grid[0])
}

type Patterns []Pattern

func ParseMirrors(r io.Reader) (Patterns, error) {
	scr := bufio.NewScanner(r)
	patterns := []Pattern{}
	curP := Pattern{grid: []string{}}
	for scr.Scan() {
		if scr.Err() != nil {
			return patterns, fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		line := scr.Text()
		if line == "" {
			patterns = append(patterns, curP)
			curP = Pattern{grid: []string{}}
			continue
		}

		curP.grid = append(curP.grid, line)
	}
	patterns = append(patterns, curP)
	return patterns, nil
}

func (p Patterns) Summarize() int {
	leftColumns := 0
	topRows := 0
	for i, pattern := range p {
		orientation, mirrorIdx := pattern.IndexMirror()
		if mirrorIdx < 0 {
			fmt.Printf("mirror not found: %d\n", i)
			continue
		}
		if orientation == Vertical {
			leftColumns += mirrorIdx
		}
		if orientation == Horizontal {
			topRows += mirrorIdx
		}
	}
	return leftColumns + (topRows * 100)
}

func (p Pattern) IndexMirror() (o Orientation, idx int) {
	if i := p.indexVerticalMirror(); i > -1 {
		return Vertical, i
	}
	return Horizontal, p.indexHorizontalMirror()
}

// IndexHorizontalMirror returns the index of the horizontal mirror, if exists.
// If a mirror does not exists, -1 is returned. The mirror will be below the returned index.
// For example, if 2 is returned, the mirror is located between rows 2 and 3.
func (p Pattern) indexHorizontalMirror() int {
	for y := 1; y < p.height(); y++ {
		if p.isMirrorAtY(y) {
			return y
		}
	}
	return -1
}

func (p Pattern) isMirrorAtY(y int) bool {
	spread := 0
	inBounds := func(spread int) bool { return y+spread < p.height() && y-spread > 0 }

	for inBounds(spread) {
		top := p.grid[y-spread-1]
		bottom := p.grid[y+spread]
		if top != bottom {
			return false
		}
		spread++
	}
	return true
}

// IndexVerticalMirror returns the index of the vertical mirror, if exists.
// If a mirror does not exists, -1 is returned. The mirror will be to the right of the returned index.
// For example, if 2 is returned, the mirror is located between columns 2 and 3.
func (p Pattern) indexVerticalMirror() int {
	for x := 1; x < p.width(); x++ {
		if p.isMirrorAtX(x) {
			return x
		}
	}
	return -1
}

func (p Pattern) isMirrorAtX(x int) bool {
	spread := 0
	inBounds := func(spread int) bool { return x+spread < p.width() && x-spread > 0 }

	for inBounds(spread) {
		left := p.columnAt(x - spread - 1)
		right := p.columnAt(x + spread)
		if left != right {
			return false
		}
		spread++
	}
	return true
}

func (p Pattern) columnAt(x int) string {
	var s strings.Builder
	for y := 0; y < p.height(); y++ {
		s.WriteByte(p.grid[y][x])
	}
	return s.String()
}
