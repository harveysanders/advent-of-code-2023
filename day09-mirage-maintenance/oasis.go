package oasis

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Report struct {
	Measurements []Measurement
}

func (r Report) Total(reverse bool) int {
	sum := 0
	for _, m := range r.Measurements {
		if reverse {
			sum += m.ExtrapolateReverse()
		} else {
			sum += m.Extrapolate()
		}
	}
	return sum
}

type Measurement struct {
	history []int
}

func (m Measurement) Extrapolate() int {
	return m.extrapolate(false)
}

func (m Measurement) ExtrapolateReverse() int {
	return m.extrapolate(true)
}

func (m Measurement) extrapolate(reverse bool) int {
	list := m.history
	allDiffs := [][]int{list}
	isStart := true
	for isStart || !allZeros(list) {
		isStart = false
		list = diffEach(list)
		allDiffs = append(allDiffs, list)
	}
	if reverse {
		return calcFirstValues(allDiffs)
	}
	return calcLastValues(allDiffs)
}

func calcFirstValues(allDiffs [][]int) int {
	values := []int{0}
	for i := len(allDiffs) - 1; i > -1; i-- {
		diffs := allDiffs[i]
		v := diffs[0] - values[len(allDiffs)-1-i]
		values = append(values, v)
	}
	return values[len(values)-1]
}

func calcLastValues(allDiffs [][]int) int {
	values := []int{0}
	for i := len(allDiffs) - 1; i > -1; i-- {
		diffs := allDiffs[i]
		v := diffs[len(diffs)-1] + values[len(allDiffs)-1-i]
		values = append(values, v)
	}
	return values[len(values)-1]
}

func allZeros(list []int) bool {
	for _, v := range list {
		if v != 0 {
			return false
		}
	}
	return true
}

func diffEach(list []int) []int {
	diffs := []int{}
	for i := 0; i < len(list)-1; i++ {
		next := list[i+1]
		cur := list[i]
		diffs = append(diffs, next-cur)
	}
	return diffs
}

func ParseReport(r io.Reader) (Report, error) {
	scr := bufio.NewScanner(r)
	report := Report{Measurements: []Measurement{}}
	for scr.Scan() {
		if scr.Err() != nil {
			return report, fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		line := scr.Text()
		raw := strings.Fields(line)
		history := make([]int, len(raw))
		for i, v := range raw {
			n, err := strconv.Atoi(v)
			if err != nil {
				return report, fmt.Errorf("strconv.Atoi(v): %w", err)
			}
			history[i] = n
		}
		report.Measurements = append(report.Measurements, Measurement{
			history: history,
		})
	}
	return report, nil
}
