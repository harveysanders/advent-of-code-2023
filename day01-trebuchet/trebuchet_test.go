package trebuchet_test

import (
	"embed"
	"io"
	"strings"
	"testing"

	trebuchet "github.com/harveysanders/advent-of-code-2023/day01-trebuchet"
	"github.com/stretchr/testify/require"
)

type comparator int

const (
	EQUAL comparator = iota
	GREATER_THAN
	LESS_THAN
)

//go:embed input/*
var inputFiles embed.FS

func TestParseCalibrationDocSample(t *testing.T) {
	t.Run("it parses the numbers from each line and sums the results", func(t *testing.T) {
		testCases := []struct {
			description string
			numberRE    string
			input       io.Reader
			wantSum     int
		}{
			{
				description: "part 1 sample input",
				numberRE:    `\d{1}`,
				input: strings.NewReader(`1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
`),
				wantSum: 142,
			},
			{
				description: "part 2 sample input",
				numberRE:    `(\d|one|two|three|four|five|six|seven|eight|nine){1}`,
				input: strings.NewReader(`two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
`),
				wantSum: 281,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				parser := trebuchet.New(tc.numberRE)
				gotSum, err := parser.ParseCalibrationDoc(tc.input)
				require.NoError(t, err)

				require.Equal(t, tc.wantSum, gotSum)
			})
		}
	})
}

func TestParseCalibrationDocFull(t *testing.T) {
	testCases := []struct {
		description string
		numberRE    string
		wantSum     int
		comparator  comparator
	}{
		{
			description: "part 1 full input",
			numberRE:    `\d{1}`,
			wantSum:     53651,
			comparator:  EQUAL,
		},
		{
			description: "part 2 full input",
			numberRE:    `(\d|one|two|three|four|five|six|seven|eight|nine){1}`,
			wantSum:     53896,
			comparator:  LESS_THAN, // 53896 is too high!
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			f, err := inputFiles.Open("input/input.txt")
			require.NoError(t, err)

			parser := trebuchet.New(tc.numberRE)
			gotSum, err := parser.ParseCalibrationDoc(f)
			require.NoError(t, err)

			switch tc.comparator {
			case EQUAL:
				require.Equal(t, tc.wantSum, gotSum)
			case GREATER_THAN:
				require.Greater(t, gotSum, tc.wantSum)
			case LESS_THAN:
				require.Less(t, gotSum, tc.wantSum)
			default:
				require.Equal(t, tc.wantSum, gotSum)
			}
		})
	}
}
