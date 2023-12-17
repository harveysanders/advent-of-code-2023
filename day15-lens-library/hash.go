package hash

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func Calculate(input string, seed int) int {
	cur := seed
	multiplier := 17
	for _, c := range input {
		cur = ((cur + int(c)) * multiplier) % 256
	}
	return cur
}

func SumInitSeq(r io.Reader) (int, error) {
	sequence, err := io.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("io.ReadAll: %w", err)
	}
	sequence = bytes.TrimSuffix(sequence, []byte("\n"))
	sum := 0
	steps := strings.Split(string(sequence), ",")
	for _, step := range steps {
		s := Calculate(step, 0)
		// fmt.Printf("%s becomes %d\n", step, s)
		sum += s
	}
	return sum, nil
}
