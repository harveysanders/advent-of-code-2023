package wasteland

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type NodeMap struct {
	LR    []string        // List of left/right ("L", "R") instructions use to move through the nodes.
	Nodes map[string]Node // Node name to node
	Start string          // Name of start node ("AAA").
	End   string          // Name of end node ("ZZZ").
}

func ParseNodeMap(r io.Reader) (NodeMap, error) {
	scr := bufio.NewScanner(r)
	nm := NodeMap{
		Nodes: make(map[string]Node),
		Start: "AAA",
		End:   "ZZZ",
	}

	isHeader := true
	for scr.Scan() {
		if scr.Err() != nil {
			return nm, fmt.Errorf("scr.Err(): %w", scr.Err())
		}

		line := scr.Text()
		if isHeader {
			isHeader = false
			nm.LR = strings.Split(strings.TrimSpace(line), "")
			continue
		}

		if line == "" {
			continue
		}

		p := strings.Split(line, " = ")
		if len(p) != 2 {
			return nm, fmt.Errorf("invalid node line: %q", line)
		}

		left, right, err := parseLRNodes(p[1])
		if err != nil {
			return nm, fmt.Errorf("parseLRNodes: %w", err)
		}
		name := p[0]
		node := Node{
			Name:  name,
			Left:  left,
			Right: right,
		}

		nm.Nodes[name] = node
	}
	return nm, nil
}

func parseLRNodes(rawLRNodes string) (left, right string, err error) {
	lr := strings.Split(strings.Trim(rawLRNodes, " ()"), ", ")
	if len(lr) != 2 {
		return "", "", fmt.Errorf("invalid nodes: %q", rawLRNodes)
	}
	return lr[0], lr[1], nil
}

// Traverse uses the left/right instructions to move from start to end, returning the number of steps taken.
func (m NodeMap) Traverse() (int, error) {
	startNode := m.Nodes[m.Start]

	steps, err := m.move(startNode, m.End, 0)
	if err != nil {
		return 0, fmt.Errorf("move: %w", err)
	}

	return steps, nil
}

func (m NodeMap) move(n Node, dest string, step int) (steps int, err error) {
	if n.Name == dest {
		return step, nil
	}

	dir := Direction(m.LR[step%len(m.LR)])
	var ok bool
	var next Node
	switch dir {
	case DirLeft:
		next, ok = m.Nodes[n.Left]
	case DirRight:
		next, ok = m.Nodes[n.Right]
	}
	if !ok {
		return step, fmt.Errorf("node not found: %+v, dir: %q", n, dir)
	}

	return m.move(next, dest, step+1)
}

type Node struct {
	Name  string // Three-letter string.
	Left  string // Name of left node.
	Right string // Name of right node.
}

type Direction string

const (
	DirLeft  Direction = "L"
	DirRight Direction = "R"
)
