package race_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	race "github.com/harveysanders/advent-of-code-2023/day06-wait-for-it"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	sample := strings.NewReader(`Time:      7  15   30
Distance:  9  40  200
`)

	races, err := race.Parse(sample, false)
	require.NoError(t, err)

	require.Len(t, races, 3)
	wantTimes := []int{7, 15, 30}
	wantDistances := []int{9, 40, 200}

	for i, r := range races {
		require.Equal(t, wantTimes[i], r.Time)
		require.Equal(t, wantDistances[i], r.Distance)
	}
}

func TestParsePart2(t *testing.T) {
	sample := strings.NewReader(`Time:      7  15   30
Distance:  9  40  200
`)

	races, err := race.Parse(sample, true)
	require.NoError(t, err)

	require.Len(t, races, 1)
	require.Equal(t, 71530, races[0].Time)
	require.Equal(t, 940200, races[0].Distance)
}

func TestWinningTimes(t *testing.T) {
	testCases := []struct {
		race      race.Race
		wantWaysN int
	}{
		{
			race: race.Race{
				Time:     7,
				Distance: 9,
			},
			wantWaysN: 4,
		},
		{
			race: race.Race{
				Time:     15,
				Distance: 40,
			},
			wantWaysN: 8,
		},
		{
			race: race.Race{
				Time:     30,
				Distance: 200,
			},
			wantWaysN: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("time: %d ms, record: %d mm", tc.race.Time, tc.race.Distance), func(t *testing.T) {
			times := tc.race.WinningTimes()
			require.Len(t, times, tc.wantWaysN)

		})
	}
}

var isCI = os.Getenv("CI") != ""

func TestDayMarginOfError(t *testing.T) {
	sample := strings.NewReader(`Time:      7  15   30
Distance:  9  40  200
`)
	fullInput, err := github.GetInputFile(6, !isCI)
	require.NoError(t, err)
	defer fullInput.Close()

	testCases := []struct {
		name            string
		input           io.ReadSeeker
		wantErrorMargin int
		part2           bool
		skip            bool
	}{
		{
			name:            "Part 1 sample",
			input:           sample,
			wantErrorMargin: 288,
		},
		{
			name:            "Part 1 full",
			input:           fullInput,
			wantErrorMargin: 4811940,
		},
		{
			name:            "Part 2 sample",
			input:           sample,
			wantErrorMargin: 71503,
			part2:           true,
		},
		{
			name:            "Part 2 full",
			input:           fullInput,
			wantErrorMargin: 30077773,
			part2:           true,
			skip:            isCI,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("Test %q skipped\n", tc.name)
			}
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)

			races, err := race.Parse(tc.input, tc.part2)
			require.NoError(t, err)
			got := races.MarginOfError()
			require.Equal(t, tc.wantErrorMargin, got)
		})
	}
}
