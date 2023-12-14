package mirror_test

import (
	"io"
	"strings"
	"testing"

	mirror "github.com/harveysanders/advent-of-code-2023/day13-point-of-incidence"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestMirror(t *testing.T) {
	input := strings.NewReader(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`)

	patterns, err := mirror.ParseMirrors(input)
	require.NoError(t, err)

	require.Len(t, patterns, 2)
}

func TestIndexMirror(t *testing.T) {
	testCases := []struct {
		name            string
		input           io.ReadSeeker
		wantOrientation mirror.Orientation
		wantIndex       int
	}{
		{
			name: "pattern 1, part 1",
			input: strings.NewReader(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.
`),
			wantOrientation: mirror.Vertical,
			wantIndex:       5,
		},
		{
			name: "pattern 2, part 1",
			input: strings.NewReader(`#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`),
			wantOrientation: mirror.Horizontal,
			wantIndex:       4,
		},
		{
			name: "pattern 3, part 1",
			input: strings.NewReader(`#########.##.##
##.....#.####.#
#.#.##...####..
...#####..##...
#.....#...##...
.##..#..#.##.#.
.##..#...#..#..
####..##.#..#.#
####..##.#..#.#
`),
			wantOrientation: mirror.Horizontal,
			wantIndex:       8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			patterns, err := mirror.ParseMirrors(tc.input)
			require.NoError(t, err)

			pattern := patterns[0]
			gotOrientation, gotIndex := pattern.IndexMirror()
			require.Equal(t, tc.wantOrientation.String(), gotOrientation.String())
			require.Equal(t, tc.wantIndex, gotIndex)
		})
	}
}

func TestSummarize(t *testing.T) {
	sample := strings.NewReader(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`)

	fullInput, err := github.GetInputFile(13, !github.IsCIEnv)
	require.NoError(t, err)
	defer fullInput.Close()

	testCases := []struct {
		name      string
		input     io.ReadSeeker
		wantTotal int
	}{
		{
			name:      "Sample part 1",
			input:     sample,
			wantTotal: 405,
		},
		{
			name:      "Full input part 1",
			input:     fullInput,
			wantTotal: 35521,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patterns, err := mirror.ParseMirrors(tc.input)
			require.NoError(t, err)

			gotTotal := patterns.Summarize()
			require.Equal(t, tc.wantTotal, gotTotal)
		})
	}
}
