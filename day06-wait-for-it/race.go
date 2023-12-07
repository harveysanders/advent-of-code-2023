package race

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Race struct {
	Time     int // Time alloted for race in milliseconds.
	Distance int // Record distance a boat has traveled for the given race time, in millimeters.
}

// WinningTimes returns a list of charge button durations needed to beat the race record distance.
func (r Race) WinningTimes() []int {
	winners := make(map[int]interface{})
	for i := 1; i <= r.Time; i++ {
		timeLeft := r.Time - i
		dist := i * timeLeft
		if dist > r.Distance {
			winners[i] = nil
		}
	}
	res := []int{}
	for i, _ := range winners {
		res = append(res, i)
	}
	return res
}

type Races []Race

func Parse(data io.Reader) (Races, error) {
	var races Races
	scr := bufio.NewScanner(data)
	for scr.Scan() {
		if scr.Err() != nil {
			return races, scr.Err()
		}

		line := scr.Text()
		if strings.Contains(line, "Time") {
			rawTimes := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, "Time:")))
			races = make([]Race, len(rawTimes))
			for i, v := range rawTimes {
				ms, err := strconv.Atoi(v)
				if err != nil {
					return races, fmt.Errorf("strconv.Atoi(): %w, val: %s", err, v)
				}
				races[i] = Race{Time: ms}
			}
			continue
		}

		if strings.Contains(line, "Distance") {
			rawDists := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, "Distance:")))
			for i, v := range rawDists {
				d, err := strconv.Atoi(v)
				if err != nil {
					return races, fmt.Errorf("strconv.Atoi(): %w, val: %s", err, v)
				}
				races[i].Distance = d
			}
			continue
		}
	}
	return races, nil
}

func (r Races) MarginOfError() int {
	res := 1
	for _, race := range r {
		res *= len(race.WinningTimes())
	}
	return res
}
