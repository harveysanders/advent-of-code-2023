package cubes_test

import (
	"testing"

	cubes "github.com/harveysanders/advent-of-code-2023/day02-cube-conundrum"
	"github.com/stretchr/testify/require"
)

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
