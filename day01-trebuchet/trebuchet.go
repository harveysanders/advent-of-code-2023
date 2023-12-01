package trebuchet

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func main() {

}

func ParseCalibrationDoc(input io.Reader) (int, error) {
	scr := bufio.NewScanner(input)

	re := regexp.MustCompile(`\d{1}`)
	var nums []int

	for scr.Scan() {
		line := scr.Bytes()
		matches := re.FindAll(line, -1)
		var digits []byte
		switch len(matches) {
		case 0:
			return 0, fmt.Errorf("no digits found in %q", string(line))
		case 1:
			digits = append(matches[0], matches[0]...)
		default:
			digits = append(matches[0], matches[len(matches)-1]...)
		}

		n, err := strconv.Atoi(string(digits))
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
