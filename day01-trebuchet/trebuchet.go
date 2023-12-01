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
	numRE *regexp.Regexp
}

func New(numberRegExp string) *Treb {
	return &Treb{
		numRE: regexp.MustCompile(numberRegExp),
	}
}

func (t Treb) ParseCalibrationDoc(input io.Reader) (int, error) {
	scr := bufio.NewScanner(input)

	var nums []int

	for scr.Scan() {
		line := scr.Bytes()
		matches := t.numRE.FindAll(line, -1)
		var digits strings.Builder
		var firstDigit string
		var lastDigit string

		switch len(matches) {
		case 0:
			return 0, fmt.Errorf("no digits found in %q", string(line))
		case 1:
			firstDigit = wordToN(string(matches[0]))
			lastDigit = firstDigit
		default:
			firstDigit = wordToN(string(matches[0]))
			lastDigit = wordToN(string(matches[len(matches)-1]))
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

// WordToN takes a spelled out number and converts it to the arabic numeral representation. If the input can not be transformed to the arabic represent, the value is returned unchanged.
// Example:
// wordToN("three") -> "3"
// wordToN("3") -> "3"
// wordToN("taco") -> "taco"
func wordToN(word string) string {
	nums := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
	}

	n, ok := nums[string(word)]
	if !ok {
		return string(word)
	}
	return n
}
