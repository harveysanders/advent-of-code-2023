package engine_test

import (
	"embed"
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

//go:embed input/*
var inputFiles embed.FS

func TestSumPartNums(t *testing.T) {
	sampleInput := io.NopCloser(strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`))
	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		input   io.ReadCloser
		wantSum int
	}{
		{
			name:    "Part 1 sample",
			input:   sampleInput,
			wantSum: 4361,
		},
		{
			name:    "Part 1 Full",
			input:   fullInput,
			wantSum: 539433,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.input.Close()

			schematics := engine.Schematic{}
			err := schematics.Parse(tc.input)
			require.NoError(t, err)

			got, err := schematics.PartNumSum()
			require.NoError(t, err)
			require.Equal(t, tc.wantSum, got)
		})
	}
}

func TestFindGears(t *testing.T) {
	sampleInput := io.NopCloser(strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`))
	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		input   io.ReadCloser
		wantSum int
	}{
		{
			name:    "Part 2 sample",
			input:   sampleInput,
			wantSum: 467835,
		},
		{
			name:    "Part 2 Full",
			input:   fullInput,
			wantSum: 75847567,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.input.Close()

			schematics := engine.Schematic{}
			err := schematics.Parse(tc.input)
			require.NoError(t, err)

			got, err := schematics.FindGears()
			require.NoError(t, err)
			require.Equal(t, tc.wantSum, got)
		})
	}
}
