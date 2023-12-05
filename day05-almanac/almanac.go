package almanac

import (
	"fmt"

	"github.com/harveysanders/advent-of-code-2023/day05-almanac/category"
)

type Almanac struct {
	Maps map[category.Name]Conversion
}

type Conversion struct {
	Src    category.Name
	Dst    category.Name
	Ranges []Range
}

type Range struct {
	SrcStart int
	DstStart int
	Length   int
}

func (c Conversion) convert(src category.Name, srcVal int) (dst category.Name, result int) {
	for _, r := range c.Ranges {
		if srcVal >= r.SrcStart && srcVal < r.SrcStart+r.Length {
			diff := r.DstStart - r.SrcStart
			return c.Dst, srcVal + diff
		}
	}
	return c.Dst, srcVal
}

func (a Almanac) ConvertTo(src category.Name, dest category.Name, val int) (int, error) {
	converter, ok := a.Maps[src]
	if !ok {
		return 0, fmt.Errorf("converter not found for %q", src)
	}
	nextDst, nextVal := converter.convert(src, val)
	for dest != category.Name(nextDst) {
		nextSrc := converter.Dst
		converter, ok = a.Maps[nextSrc]
		if !ok {
			return 0, fmt.Errorf("converter not found for %q", nextSrc)
		}
		nextDst, nextVal = converter.convert(converter.Src, nextVal)
	}
	return nextVal, nil
}
