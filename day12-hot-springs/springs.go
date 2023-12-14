package springs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Record struct {
	// Row of springs, identified by their known state.
	//	"#" - damaged
	//	"." - operational
	//	"?" - unknown
	Conditions []string
	// Ordered list of the size of each contiguous group of damaged springs.
	DamagedGroupSizes []int
}

func parseRecord(row string) (Record, error) {
	p := strings.Fields(row)
	if len(p) < 2 {
		return Record{}, fmt.Errorf("invalid row: %q", row)
	}
	rawGroups := strings.Split(p[1], ",")
	groups := make([]int, len(rawGroups))
	for i, v := range rawGroups {
		n, err := strconv.Atoi(v)
		if err != nil {
			return Record{}, fmt.Errorf("strconv.Atoi(v): %w", err)
		}
		groups[i] = n
	}
	r := Record{
		Conditions:        strings.Split(p[0], ""),
		DamagedGroupSizes: groups,
	}
	return r, nil
}

type Records []Record

func ParseRecords(r io.Reader) (Records, error) {
	scr := bufio.NewScanner(r)
	records := []Record{}
	for scr.Scan() {
		if scr.Err() != nil {
			return records, fmt.Errorf("scr.Err(): %w", scr.Err())
		}
		row := scr.Text()
		record, err := parseRecord(row)
		if err != nil {
			return records, fmt.Errorf("parseRecord: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}
