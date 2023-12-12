package maze_test

import (
	"io"
	"strings"
	"testing"

	maze "github.com/harveysanders/advent-of-code-2023/day10-pipe-maze"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestFindStart(t *testing.T) {
	sample1 := strings.NewReader(`-L|F7
7S-7|
L|7||
-L-J|
L|-JF
`)
	sample2 := strings.NewReader(`7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ
`)

	testCases := []struct {
		name      string
		input     io.Reader
		wantStart maze.Coord
	}{
		{
			name:      "sample 1",
			input:     sample1,
			wantStart: maze.Coord{X: 1, Y: 1},
		},
		{
			name:      "sample 2",
			input:     sample2,
			wantStart: maze.Coord{X: 0, Y: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := maze.ParseMaze(tc.input)
			require.NoError(t, err)

			gotStart, err := m.FindStart()
			require.NoError(t, err)
			require.Equal(t, tc.wantStart, gotStart)
		})
	}
}

func TestFarthestDist(t *testing.T) {
	sample1 := strings.NewReader(`.....
.S-7.
.|.|.
.L-J.
.....
`)
	sample2 := strings.NewReader(`..F7.
.FJ|.
SJ.L7
|F--J
LJ...
`)

	fullInput, err := github.GetInputFile(10, !github.IsCIEnv)
	require.NoError(t, err)

	defer fullInput.Close()

	testCases := []struct {
		name     string
		input    io.ReadSeeker
		wantDist int
	}{
		{
			name:     "sample 1, part 1",
			input:    sample1,
			wantDist: 4,
		},
		{
			name:     "sample 2, part 1",
			input:    sample2,
			wantDist: 8,
		},
		{
			name:     "full input, part 1",
			input:    fullInput,
			wantDist: 6856,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)

			m, err := maze.ParseMaze(tc.input)
			require.NoError(t, err)

			gotDist, err := m.FarthestDistFromStart()
			require.NoError(t, err)

			require.Equal(t, tc.wantDist, gotDist)
		})
	}
}
