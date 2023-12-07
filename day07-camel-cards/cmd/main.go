package main

import (
	"log"

	camel "github.com/harveysanders/advent-of-code-2023/day07-camel-cards"
	"github.com/harveysanders/advent-of-code-2023/internal/github"
)

func main() {
	fullInput, err := github.GetInputFile(7, false)
	if err != nil {
		log.Fatal(err)
	}
	defer fullInput.Close()

	hands, err := camel.ParseHands(fullInput)
	if err != nil {
		log.Fatal(err)
	}

	w := hands.TotalWinnings()
	log.Printf("TOTAL: %d", w)
}
