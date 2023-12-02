package trebuchet

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Treb struct {
	includeWords bool // If true, spelled out numbers (i.e. "eight") will count as valid digits.
}

// New creates a new calibration document parser. Passing true for includeWords will inform the parser to include spelled-out numbers while parsing.
func New(includeWords bool) *Treb {
	return &Treb{
		includeWords: includeWords,
	}
}

func (t Treb) ParseCalibrationDoc(input io.Reader) (int, error) {
	scr := bufio.NewScanner(input)
	digitRE := regexp.MustCompile(`\d{1}`)
	// Go's RegExp implementation does not support lookahead, I'm using this
	// markers map to "mark" the English representation with the Arabic numeral.
	// Since the arabic numeral are near the middle of the words, this will handle overlaps like "oneight" in "bhdf315nineoneightzlp".
	// "oneight" -> "o1ei8ht" -> [["1"],["8"]]
	markers := map[string]string{
		"one":   "o1e",
		"two":   "t2o",
		"three": "th3ee",
		"four":  "f4ur",
		"five":  "f5ve",
		"six":   "s6x",
		"seven": "se7en",
		"eight": "ei8ht",
		"nine":  "n9ne",
	}

	var nums []int
	for scr.Scan() {
		line := scr.Text()

		if t.includeWords {
			for old, new := range markers {
				line = strings.ReplaceAll(line, old, new)
			}
		}

		matches := digitRE.FindAllString(line, -1)
		var digits strings.Builder
		var firstDigit string
		var lastDigit string

		switch len(matches) {
		case 0:
			return 0, fmt.Errorf("no digits found in %q", line)
		case 1:
			firstDigit = matches[0]
			lastDigit = firstDigit
		default:
			firstDigit = matches[0]
			lastDigit = matches[len(matches)-1]
		}

		digits.WriteString(firstDigit)
		digits.WriteString(lastDigit)

		n, err := strconv.Atoi(digits.String())
		if err != nil {
			return 0, fmt.Errorf("strconv.Atoi: %w", err)
		}
		nums = append(nums, n)
	}

	sum := 0
	for _, v := range nums {
		sum += v
	}

	return sum, nil
}
