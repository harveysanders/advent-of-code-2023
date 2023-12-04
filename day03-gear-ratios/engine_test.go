package engine_test

import (
	"io"
	"strings"
	"testing"

	engine "github.com/harveysanders/advent-of-code-2023/day03-gear-ratios"
	"github.com/stretchr/testify/require"
)

func TestCollectNumbers(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`))

	schematics := engine.Schematic{}

	err := schematics.Parse(input)
	require.NoError(t, err)

	nums, err := schematics.CollectNumbers()
	require.NoError(t, err)

	require.Len(t, nums, 10)
	thirtyFive := nums[2]
	require.Equal(t, 35, thirtyFive.Value)
	wantLoc := engine.Coord{2, 2}
	require.Equal(t, wantLoc, thirtyFive.Location)
	require.Equal(t, 2, thirtyFive.Size)
}

func TestIsPartNum(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`))

	schematics := engine.Schematic{}

	err := schematics.Parse(input)
	require.NoError(t, err)

	nums, err := schematics.CollectNumbers()
	thirtyFive := nums[2]
	require.NoError(t, err)
	isPartNum := schematics.IsPartNum(thirtyFive)

	require.True(t, isPartNum)
}

func TestSumPartNums(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`))

	schematics := engine.Schematic{}
	err := schematics.Parse(input)
	require.NoError(t, err)

	wantPartNums := []int{467, 35, 633, 617, 592, 755, 664, 598}
	got, err := schematics.PartNums()
	require.NoError(t, err)

	require.Equal(t, wantPartNums, got)
}
