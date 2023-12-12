package oasis_test

import (
	"io"
	"strings"
	"testing"

	oasis "github.com/harveysanders/advent-of-code-2023/day09-mirage-maintenance"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestExtrapolate(t *testing.T) {
	input := strings.NewReader(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`)
	report, err := oasis.ParseReport(input)
	require.NoError(t, err)

	require.Len(t, report.Measurements, 3)
	wantExtrapolatedValues := []int{18, 28, 68}

	for i, wantV := range wantExtrapolatedValues {
		got := report.Measurements[i].Extrapolate()
		require.Equal(t, wantV, got)
	}
}

func TestExtrapolateReverse(t *testing.T) {
	input := strings.NewReader(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`)
	report, err := oasis.ParseReport(input)
	require.NoError(t, err)

	wantExtrapolatedValues := []int{-3, 0, 5}
	for i, wantV := range wantExtrapolatedValues {
		got := report.Measurements[i].ExtrapolateReverse()
		require.Equal(t, wantV, got)
	}
}

func TestTotal(t *testing.T) {
	sample := strings.NewReader(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`)

	fullInput, err := github.GetInputFile(9, !github.IsCIEnv)
	require.NoError(t, err)
	defer fullInput.Close()

	testCases := []struct {
		name                 string
		input                io.ReadSeeker
		wantTotal            int
		reverseExtrapolation bool
	}{
		{
			name:      "sample part 1",
			input:     sample,
			wantTotal: 114,
		},
		{
			name:      "full input part 1",
			input:     fullInput,
			wantTotal: 2075724761,
		},
		{
			name:                 "sample part 2",
			input:                sample,
			wantTotal:            2,
			reverseExtrapolation: true,
		},
		{
			name:                 "full input part 2",
			input:                fullInput,
			wantTotal:            1072,
			reverseExtrapolation: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)

			report, err := oasis.ParseReport(tc.input)
			require.NoError(t, err)

			got := report.Total(tc.reverseExtrapolation)
			require.Equal(t, tc.wantTotal, got)
		})
	}
}
