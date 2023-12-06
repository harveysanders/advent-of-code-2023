package main

import (
	"fmt"
	"log"

	almanac "github.com/harveysanders/advent-of-code-2023/day05-almanac"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
)

func main() {
	input, err := github.GetInputFile(5, true)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()

	a, err := almanac.Parse(input)
	if err != nil {
		log.Fatal(err)
	}

	lowest, err := a.LowestLocation(true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lowest)
}
