package cubes_test

import (
	"embed"
	"io"
	"strings"
	"testing"

	cubes "github.com/harveysanders/advent-of-code-2023/day02-cube-conundrum"
	"github.com/stretchr/testify/require"
)

//go:embed input/*.txt
var inputFiles embed.FS

func TestValidateGame(t *testing.T) {
	testCases := []struct {
		bag          *cubes.Bag
		game         cubes.Game
		wantPossible bool
	}{
		{
			bag: cubes.NewBag(12, 13, 14),
			game: cubes.Game{
				Sets: []cubes.Set{{Red: 20}},
			},
		},
	}

	for _, tc := range testCases {
		gotPossible := tc.bag.ValidateGame(tc.game)
		require.Equal(t, tc.wantPossible, gotPossible)
	}
}

func TestGameUnmarshal(t *testing.T) {
	testCases := []struct {
		rawGame  string
		wantGame cubes.Game
	}{
		{
			rawGame: "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
			wantGame: cubes.Game{
				ID: 1,
				Sets: []cubes.Set{
					{Red: 4, Blue: 3},
					{Red: 1, Green: 2, Blue: 6},
					{Green: 2},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.rawGame, func(t *testing.T) {
			g := cubes.Game{}

			err := g.Parse(tc.rawGame)
			require.NoError(t, err)
			require.Equal(t, tc.wantGame.ID, g.ID, "wrong ID")
			require.Len(t, g.Sets, len(tc.wantGame.Sets))
			for i, wantSet := range tc.wantGame.Sets {
				require.Equal(t, wantSet, g.Sets[i])
			}
		})
	}
}

func TestSampleRecordDecode(t *testing.T) {
	testCases := []struct {
		desc    string
		input   io.Reader
		bag     cubes.Bag
		wantSum int
	}{
		{
			desc: "Part 1 sample",
			input: strings.NewReader(`Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
`),
			bag:     *cubes.NewBag(12, 13, 14),
			wantSum: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var record cubes.Record
			err := record.Decode(tc.input)
			require.NoError(t, err)

			ids := record.ValidGameIDs(tc.bag)
			sum := 0
			for _, id := range ids {
				sum += id
			}

			require.Equal(t, tc.wantSum, sum)
		})
	}
}

func TestPart1FullInput(t *testing.T) {
	t.Run("part 1 full input", func(t *testing.T) {
		input, err := inputFiles.Open("input/input.txt")
		require.NoError(t, err)

		var record cubes.Record
		err = record.Decode(input)
		require.NoError(t, err)

		ids := record.ValidGameIDs(*cubes.NewBag(12, 13, 14))

		sum := 0
		for _, id := range ids {
			sum += id
		}

		require.Equal(t, 2795, sum)
	})
}

func TestPart2Sample(t *testing.T) {
	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		wantSum int
		desc    string
		input   io.ReadCloser
	}{
		{
			desc:    "part 2 sample",
			wantSum: 2286,
			input: io.NopCloser(strings.NewReader(`Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
`)),
		},
		{
			desc:    "part 2 full input",
			wantSum: 75561,
			input:   fullInput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer tc.input.Close()

			var record cubes.Record
			err := record.Decode(tc.input)
			require.NoError(t, err)

			sum := record.Part2()
			require.Equal(t, tc.wantSum, sum)
		})
	}

}
