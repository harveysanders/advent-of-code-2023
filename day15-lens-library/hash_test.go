package hash_test

import (
	"io"
	"strings"
	"testing"

	hash "github.com/harveysanders/advent-of-code-2023/day15-lens-library"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {
	input := `HASH`

	got := hash.Calculate(input, 0)
	require.Equal(t, 52, got)
}

func TestSumInitSeq(t *testing.T) {
	sample := strings.NewReader(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7
`)
	fullInput, err := github.GetInputFile(15, !github.IsCIEnv)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		input   io.ReadSeeker
		wantSum int
	}{
		{
			name:    "sample, part 1",
			input:   sample,
			wantSum: 1320,
		},
		{
			name:    "full, part 1",
			input:   fullInput,
			wantSum: 495972,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)
			got, err := hash.SumInitSeq(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.wantSum, got)
		})
	}
}
