package almanac_test

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"

	almanac "github.com/harveysanders/advent-of-code-2023/day05-almanac"
	category "github.com/harveysanders/advent-of-code-2023/day05-almanac/category"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		src      category.Name
		dst      category.Name
		startVal int
		wantVal  int
	}{
		{
			src:      category.Seed,
			startVal: 79,
			dst:      category.Soil,
			wantVal:  81,
		},
		{
			src:      category.Seed,
			startVal: 14,
			dst:      category.Soil,
			wantVal:  14,
		},
		{
			src:      category.Seed,
			startVal: 55,
			dst:      category.Soil,
			wantVal:  57,
		},
		{
			src:      category.Seed,
			startVal: 13,
			dst:      category.Soil,
			wantVal:  13,
		},
		{
			src:      category.Seed,
			startVal: 79,
			dst:      category.Fertilizer,
			wantVal:  81,
		},
	}

	a := almanac.Almanac{
		Maps: map[category.Name]almanac.Conversion{
			category.Seed: {
				Src: category.Seed,
				Dst: category.Soil,
				Ranges: []almanac.Range{
					{DstStart: 50, SrcStart: 98, Length: 2},
					{DstStart: 52, SrcStart: 50, Length: 48},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("src: %q, dest: %q, value: %d", tc.src, tc.dst, tc.startVal), func(t *testing.T) {
			gotVal, err := a.ConvertTo(tc.src, tc.dst, tc.startVal)
			require.NoError(t, err)
			require.Equal(t, tc.wantVal, gotVal)
		})
	}
}

func TestParse(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`))

	a, err := almanac.Parse(input)
	require.NoError(t, err)

	wantSeeds := []int{79, 14, 55, 13}
	require.Equal(t, wantSeeds, a.Seeds)

	require.Len(t, a.Maps, 7)
	sources := []category.Name{
		category.Seed,
		category.Soil,
		category.Fertilizer,
		category.Water,
		category.Light,
		category.Temperature,
		category.Humidity,
	}

	for _, s := range sources {
		_, ok := a.Maps[s]
		require.Truef(t, ok, "%s", s)
	}
}

func TestCovertWithSample(t *testing.T) {
	input := io.NopCloser(strings.NewReader(`seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`))

	testCases := []struct {
		src      category.Name
		dst      category.Name
		startVal int
		wantVal  int
	}{
		{
			src:      category.Seed,
			dst:      category.Location,
			startVal: 79,
			wantVal:  82,
		},
	}

	a, err := almanac.Parse(input)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("src: %q, dest: %q, value: %d", tc.src, tc.dst, tc.startVal), func(t *testing.T) {
			gotVal, err := a.ConvertTo(tc.src, tc.dst, tc.startVal)
			require.NoError(t, err)
			require.Equal(t, tc.wantVal, gotVal)
		})
	}
}

//go:embed input/*.txt
var inputFiles embed.FS

func TestLowest(t *testing.T) {
	sample := io.NopCloser(strings.NewReader(`seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`))

	fullInput, err := inputFiles.Open("input/input.txt")
	require.NoError(t, err)

	testCases := []struct {
		name  string
		want  int
		input io.ReadCloser
	}{
		{
			name:  "sample part 1",
			want:  35,
			input: sample,
		},
		{
			name:  "full part 1",
			want:  88151870,
			input: fullInput,
		},
	}

	for _, tc := range testCases {
		defer tc.input.Close()

		a, err := almanac.Parse(tc.input)
		require.NoError(t, err)

		got, err := a.LowestLocation()
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}
