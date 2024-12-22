package main

import (
	"fmt"

	"github.com/llravell/simple-cards/pkg/quizlet"
)

func main() {
	parser, err := quizlet.NewParser()
	if err != nil {
		panic(err)
	}

	cards, err := parser.Parse("768736583")
	if err != nil {
		panic(err)
	}

	fmt.Println(cards)
}
