package trebuchet_test

import (
	"embed"
	"io"
	"strings"
	"testing"

	trebuchet "github.com/harveysanders/advent-of-code-2023/day01-trebuchet"
	"github.com/stretchr/testify/require"
)

//go:embed input/*
var inputFiles embed.FS

func TestParseCalibrationDoc(t *testing.T) {
	t.Run("it parses the numbers from each line and sums the results", func(t *testing.T) {
		testCases := []struct {
			description string
			input       io.Reader
			wantSum     int
		}{
			{
				description: "sample input",
				input: strings.NewReader(`1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
`),
				wantSum: 142,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				gotSum, err := trebuchet.ParseCalibrationDoc(tc.input)
				require.NoError(t, err)

				require.Equal(t, tc.wantSum, gotSum)
			})
		}
	})

	t.Run("works on part1 input file", func(t *testing.T) {
		f, err := inputFiles.Open("input/part1.txt")
		require.NoError(t, err)

		gotSum, err := trebuchet.ParseCalibrationDoc(f)
		require.NoError(t, err)

		require.Equal(t, 53651, gotSum)
	})
}
