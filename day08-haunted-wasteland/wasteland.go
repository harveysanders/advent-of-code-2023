package wasteland

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

type NodeMap struct {
	LR    []string        // List of left/right ("L", "R") instructions use to move through the nodes.
	Nodes map[string]Node // Node name to node

	lock   *sync.Mutex
	ghosts []ghost
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

type ghost struct {
	curStep int
	curNode Node
}

func ParseNodeMap(r io.Reader) (NodeMap, error) {
	scr := bufio.NewScanner(r)
	nm := NodeMap{
		lock:  &sync.Mutex{},
		Nodes: make(map[string]Node),
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

// TraverseSingle uses the left/right instructions to move from start to end, returning the number of steps taken.
func (m NodeMap) TraverseSingle(start, end string) (int, error) {
	startNode := m.Nodes[start]

	steps, err := m.move(startNode, end, 0)
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

func (m *NodeMap) moveGhost(idx int) {
	g := m.ghosts[idx]
	dir := Direction(m.LR[g.curStep%len(m.LR)])
	var ok bool
	var next Node
	switch dir {
	case DirLeft:
		next, ok = m.Nodes[g.curNode.Left]
	case DirRight:
		next, ok = m.Nodes[g.curNode.Right]
	}
	if !ok {
		log.Fatalf("node not found: %+v, dir: %q", g.curNode, dir)
	}
	g.curNode = next
	g.curStep++
	m.lock.Lock()
	m.ghosts[idx] = g
	m.lock.Unlock()

	if strings.HasSuffix(next.Name, "Z") {
		log.Printf("ghost %d at end %q\n", idx, next.Name)
	}
}

func (n *NodeMap) TraverseParallel(start, end string) (int, error) {
	n.ghosts = make([]ghost, 0)
	for name, node := range n.Nodes {
		if strings.HasSuffix(name, start) {
			n.ghosts = append(n.ghosts, ghost{curNode: node})
		}
	}

	// Spin up a go routine for each of the start nodes
	// run them each one step at a time until they all are on a node that ends with the end parameter ("Z")
	for {
		var wg sync.WaitGroup
		for i, _ := range n.ghosts {
			wg.Add(1)
			go func(n *NodeMap, idx int) {
				n.moveGhost(idx)
				wg.Done()
			}(n, i)
		}

		wg.Wait()

		allDone := true
		for i, v := range n.ghosts {
			if !strings.HasSuffix(v.curNode.Name, "Z") {
				allDone = false
				break
			}
			if i > 0 {
				log.Printf("ghost %d of %d, step: %d, node: %q\n", i, len(n.ghosts), v.curStep, v.curNode.Name)
			}
		}
		if allDone {
			return n.ghosts[0].curStep, nil
		}
	}

}
