package maze

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
)

type Maze struct {
	grid   []string
	height int
	width  int
	route  []Pipe
}

type Coord struct {
	X int
	Y int
}

type Pipe struct {
	con        Connector
	directions []Direction
	loc        Coord
}

func NewPipe(label string, x, y int) Pipe {
	conDirections := map[Connector][]Direction{
		ConnVertical:   {North, South},
		ConnHorizontal: {East, West},
		ConnL:          {North, East},
		ConnJ:          {North, West},
		Conn7:          {South, West},
		ConnF:          {East, South},
		ConnStart:      {North, East, South, West},
	}
	p := Pipe{
		con: Connector(label),
		loc: Coord{X: x, Y: y},
	}
	p.directions = conDirections[p.con]
	return p
}

type Connector string

const (
	ConnVertical   Connector = "|" // "|" connects north and south.
	ConnHorizontal Connector = "-" // "-" connects east and west.
	ConnL          Connector = "L" // "L" is a 90-degree bend connects north and east.
	ConnJ          Connector = "J" // "J" is a 90-degree bend connects north and west.
	Conn7          Connector = "7" // "7" is a 90-degree bend connects south and west.
	ConnF          Connector = "F" // "F" is a 90-degree bend connects south and east.
	ConnStart      Connector = "S" // "S" represents the start of the maze and can be any connector.
)

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func ParseMaze(r io.Reader) (Maze, error) {
	maze := Maze{
		grid:  make([]string, 0),
		route: make([]Pipe, 0),
	}
	scr := bufio.NewScanner(r)
	for scr.Scan() {
		if scr.Err() != nil {
			return maze, fmt.Errorf("scr.Err: %w", scr.Err())
		}

		line := scr.Text()
		maze.grid = append(maze.grid, line)
	}

	if len(maze.grid) > 0 {
		maze.width = len(maze.grid[0])
	}
	maze.height = len(maze.grid)
	return maze, nil
}

func (m Maze) FindStart() (Coord, error) {
	startChar := "S"
	for y, line := range m.grid {
		x := strings.Index(line, startChar)
		if x >= 0 {
			return Coord{X: x, Y: y}, nil
		}
	}
	return Coord{}, fmt.Errorf("start %q not found", startChar)
}

func (m Maze) FarthestDistFromStart() (int, error) {
	startLoc, err := m.FindStart()
	if err != nil {
		return 0, fmt.Errorf("m.FindStart(): %w", err)
	}

	startLabel := string(m.grid[startLoc.Y][startLoc.X])
	start := NewPipe(startLabel, startLoc.X, startLoc.Y)
	curPipe := start
	isStart := true
	var fromDir Direction
	for isStart || startLabel != string(curPipe.con) {
		isStart = false
		for _, nextDir := range curPipe.directions {
			if nextDir == fromDir {
				// Keep looking if we're facing the direction we just came from.
				continue
			}
			nextVal, nextPos, ok := m.Move(curPipe.loc, nextDir)
			if ok {
				if nextVal == "." {
					continue
				}
				next := NewPipe(nextVal, nextPos.X, nextPos.Y)
				if next.connects(nextDir) {
					if next.con != ConnStart {
						m.route = append(m.route, next)
					}
					curPipe = next
					// Where we came from, ex: if we just moved to the the east (nextDir), we came from the west.
					fromDir = nextDir.complement()
					break
				}
			}
		}
	}

	dist := 0
	if len(m.route) > 0 {
		// Include +1 back to start step
		dist = int(math.Ceil(float64(len(m.route)+1) / 2))
	}
	return dist, nil
}

func (m Maze) Move(loc Coord, dir Direction) (val string, next Coord, ok bool) {
	dist := 1
	switch dir {
	case North:
		return m.north(loc, dist)
	case South:
		return m.south(loc, dist)
	case West:
		return m.west(loc, dist)
	case East:
		return m.east(loc, dist)
	}
	return val, next, false
}

func (m Maze) north(loc Coord, dist int) (val string, next Coord, ok bool) {
	nextY := loc.Y - dist
	if nextY < 0 {
		return val, next, false
	}
	return string(m.grid[nextY][loc.X]), Coord{X: loc.X, Y: nextY}, true
}

func (m Maze) south(loc Coord, dist int) (val string, next Coord, ok bool) {
	nextY := loc.Y + dist
	if nextY > m.height-1 {
		return val, next, false
	}
	return string(m.grid[nextY][loc.X]), Coord{X: loc.X, Y: nextY}, true
}

func (m Maze) east(loc Coord, dist int) (val string, next Coord, ok bool) {
	nextX := loc.X + dist
	if nextX > m.width-1 {
		return val, next, false
	}
	return string(m.grid[loc.Y][nextX]), Coord{X: nextX, Y: loc.Y}, true
}

func (m Maze) west(loc Coord, dist int) (val string, next Coord, ok bool) {
	nextX := loc.X - dist
	if nextX < 0 {
		return val, next, false
	}
	return string(m.grid[loc.Y][nextX]), Coord{X: nextX, Y: loc.Y}, true
}

func (p Pipe) connects(from Direction) bool {
	oppDirection := from.complement()
	return slices.Contains(p.directions, oppDirection)
}

func (d Direction) complement() Direction {
	c := map[Direction]Direction{
		North: South,
		South: North,
		West:  East,
		East:  West,
	}

	return c[d]
}
