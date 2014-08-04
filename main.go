package main

import (
	"fmt"
)

func main() {
	experiment := NewExperiment("console")
	experiment.AddBandit(NewBandit("control-pre-avatar", 15076, 481687))
	experiment.AddBandit(NewBandit("test-pre-avatar", 5011, 161368))

	winner := pickOptimalVariant(experiment, 1000)
	fmt.Printf("%v\n", winner)
}
