package hash

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type HASHMAP struct {
	boxes [256]box
}

func New() HASHMAP {
	hm := &HASHMAP{
		boxes: [256]box{},
	}
	return *hm
}

type box []lens

type lens struct {
	label    string
	focalLen string
}

func (h *HASHMAP) Initialize(r io.Reader) error {
	sequence, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	sequence = bytes.TrimSuffix(sequence, []byte("\n"))
	steps := strings.Split(string(sequence), ",")
	for _, step := range steps {
		h.Apply(step)
	}
	return nil
}

func (h *HASHMAP) Apply(step string) {
	label, focalLen, isSetOp := strings.Cut(step, "=")
	// Check if set operation
	if isSetOp {
		boxID := Calculate(label)
		lens := lens{label: label, focalLen: focalLen}
		box := h.boxes[boxID]
		idx := slices.IndexFunc(box, containsLabel(label))
		if idx == -1 {
			box = append(box, lens)
		} else {
			box[idx] = lens
		}
		h.boxes[boxID] = box
		return
	}

	// Check if remove operation
	label, _, isRemoveOp := strings.Cut(step, "-")
	if isRemoveOp {
		boxID := Calculate(label)
		box := h.boxes[boxID]
		idx := slices.IndexFunc(box, containsLabel(label))
		if idx == -1 {
			return
		}
		box = slices.Delete(box, idx, idx+1)
		h.boxes[boxID] = box
	}
}

func containsLabel(label string) func(lens) bool {
	return func(l lens) bool { return l.label == label }
}

func (h *HASHMAP) FocusingPower() (int, error) {
	sum := 0
	for boxID, b := range h.boxes {
		if len(b) == 0 {
			continue
		}
		for i, lens := range b {
			// The focal length of the lens.
			focalLen, err := strconv.Atoi(lens.focalLen)
			if err != nil {
				return 0, fmt.Errorf("strconv.Atoi: %w", err)
			}
			// The slot number of the lens within the box: 1 for the first lens, 2 for the second lens, and so on.
			slotID := i + 1
			// One plus the box number of the lens in question.
			focusingPow := (1 + boxID) * slotID * focalLen
			sum += focusingPow
		}
	}
	return sum, nil
}

func Calculate(input string) int {
	cur := 0
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
		s := Calculate(step)
		sum += s
	}
	return sum, nil
}
