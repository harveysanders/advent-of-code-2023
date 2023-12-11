package main

import (
	"log"
	"os"

	wl "github.com/harveysanders/advent-of-code-2023/day08-haunted-wasteland"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
)

func main() {
	useLocal := os.Getenv("CI") == ""
	fullInput, err := github.GetInputFile(8, useLocal)
	if err != nil {
		log.Fatal(err)
	}
	defer fullInput.Close()

	nodeMap, err := wl.ParseNodeMap(fullInput)
	if err != nil {
		log.Fatal(err)
	}

	gotSteps, err := nodeMap.TraverseParallel("A", "Z")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("FINAL: %d\n", gotSteps)
}
