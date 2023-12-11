package wasteland_test

import (
	"io"
	"os"
	"strings"
	"testing"

	wl "github.com/harveysanders/advent-of-code-2023/day08-haunted-wasteland"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
	"github.com/stretchr/testify/require"
)

func TestTraverse(t *testing.T) {
	sample1 := strings.NewReader(`RL

AAA = (BBB, CCC)
BBB = (DDD, EEE)
CCC = (ZZZ, GGG)
DDD = (DDD, DDD)
EEE = (EEE, EEE)
GGG = (GGG, GGG)
ZZZ = (ZZZ, ZZZ)
`)

	sample2 := strings.NewReader(`LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)`)

	useLocal := os.Getenv("CI") == ""
	fullInput, err := github.GetInputFile(8, useLocal)
	require.NoError(t, err)
	defer fullInput.Close()

	testCases := []struct {
		name         string
		input        io.ReadSeeker
		wantLR       []string
		wantNodesLen int
		wantSteps    int
	}{
		{
			name:         "sample 1",
			input:        sample1,
			wantLR:       []string{"R", "L"},
			wantNodesLen: 7,
			wantSteps:    2,
		},
		{
			name:         "sample 2",
			input:        sample2,
			wantLR:       []string{"L", "L", "R"},
			wantNodesLen: 3,
			wantSteps:    6,
		},
		{
			name:         "full part 1",
			input:        fullInput,
			wantNodesLen: 766,
			wantSteps:    19951,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)

			nodeMap, err := wl.ParseNodeMap(tc.input)
			require.NoError(t, err)

			if tc.wantLR != nil {
				require.Equal(t, tc.wantLR, nodeMap.LR)
			}
			require.Len(t, nodeMap.Nodes, tc.wantNodesLen)

			gotSteps, err := nodeMap.TraverseSingle("AAA", "ZZZ")
			require.NoError(t, err)
			require.Equal(t, tc.wantSteps, gotSteps)
		})
	}
}

func TestTraverseParallel(t *testing.T) {
	sample1 := strings.NewReader(`LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)
`)

	useLocal := os.Getenv("CI") == ""
	fullInput, err := github.GetInputFile(8, useLocal)
	require.NoError(t, err)

	defer fullInput.Close()

	testCases := []struct {
		name      string
		input     io.ReadSeeker
		wantSteps int
	}{
		{
			name:      "sample 1, part 2",
			input:     sample1,
			wantSteps: 6,
		},
		{
			name:      "full input, part 2",
			input:     fullInput,
			wantSteps: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Seek(0, io.SeekStart)
			require.NoError(t, err)

			nodeMap, err := wl.ParseNodeMap(tc.input)
			require.NoError(t, err)

			gotSteps, err := nodeMap.TraverseParallel("A", "Z")
			require.NoError(t, err)
			require.Equal(t, tc.wantSteps, gotSteps)
		})

	}
}
