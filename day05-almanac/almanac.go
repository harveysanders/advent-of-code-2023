package almanac

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/harveysanders/advent-of-code-2023/day05-almanac/category"
)

type Almanac struct {
	Seeds []int
	Maps  map[category.Name]Conversion
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

func (a Almanac) LowestLocation(useRange bool) (int, error) {
	lowest := math.MaxFloat64
	if !useRange {
		for _, seed := range a.Seeds {
			location, err := a.ConvertTo(category.Seed, category.Location, seed)
			if err != nil {
				return 0, err
			}
			lowest = math.Min(lowest, float64(location))
		}
		return int(lowest), nil
	}

	// Part 2 range mode
	for i := 0; i < len(a.Seeds); i += 2 {
		startTime := time.Now()

		start := a.Seeds[i]
		count := a.Seeds[i+1]
		for seed := start; seed < start+count; seed++ {
			if seed%int(math.Pow10(6)) == 0 {
				fmt.Println(".")
			}
			location, err := a.ConvertTo(category.Seed, category.Location, seed)
			if err != nil {
				return 0, err
			}
			lowest = math.Min(lowest, float64(location))
		}

		fmt.Printf("loop: %d, elapsed(s): %d, lowest: %d\n",
			i,
			int(time.Since(startTime).Seconds()),
			int(lowest))
	}
	return int(lowest), nil
}

func Parse(r io.Reader) (Almanac, error) {
	a := Almanac{
		Seeds: make([]int, 0),
		Maps:  make(map[category.Name]Conversion),
	}
	scr := bufio.NewScanner(r)
	isHeader := true
	for scr.Scan() {
		if scr.Err() != nil {
			return a, scr.Err()
		}

		line := scr.Text()

		if isHeader {
			isHeader = false
			parts := strings.TrimPrefix(line, "seeds: ")
			nums := strings.Fields(parts)
			for _, v := range nums {
				n, err := strconv.Atoi(v)
				if err != nil {
					return a, fmt.Errorf("convert num: %w, val: %q", err, v)
				}
				a.Seeds = append(a.Seeds, n)
			}

			// Skip next empty line
			scr.Scan()
			continue
		}

		// Parse conversion map
		c, err := parseConversionMap(scr)
		if err != nil {
			return a, fmt.Errorf("parseConversionMap: %w", err)
		}
		a.Maps[c.Src] = c
	}
	return a, nil
}

func parseConversionMap(scr *bufio.Scanner) (Conversion, error) {
	c := Conversion{
		Ranges: make([]Range, 0),
	}

	line := scr.Text()
	header := strings.TrimSuffix(line, " map:")
	cats := strings.Split(strings.TrimSpace(header), "-")
	if len(cats) != 3 {
		return c, fmt.Errorf("invalid categories: %q", header)
	}
	c.Src = category.Name(cats[0])
	c.Dst = category.Name(cats[2])

	for scr.Scan() {
		if scr.Err() != nil {
			return c, scr.Err()
		}
		line := scr.Text()
		if line == "" {
			// Stop parsing
			return c, nil
		}

		maps := strings.Fields(line)
		if len(maps) != 3 {
			return c, fmt.Errorf("invalid mapping: %q", line)
		}
		rangeVals := make([]int, 3)
		for i, v := range maps {
			n, err := strconv.Atoi(v)
			if err != nil {
				return c, fmt.Errorf("strconv.Atoi: %w, val: %q", err, v)
			}
			rangeVals[i] = n
		}
		r := Range{
			DstStart: rangeVals[0],
			SrcStart: rangeVals[1],
			Length:   rangeVals[2],
		}

		c.Ranges = append(c.Ranges, r)
	}
	return c, nil
}
