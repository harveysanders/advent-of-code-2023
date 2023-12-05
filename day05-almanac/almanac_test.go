package almanac_test

import (
	"fmt"
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
