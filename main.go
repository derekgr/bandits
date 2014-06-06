package main

import (
	"fmt"
)

func main() {
	experiment := NewExperiment("google example")
	experiment.AddBandit(NewBandit("4%", 0.040))
	experiment.AddBandit(NewBandit("5%", 0.050))
	experiment.AddBandit(NewBandit("4.5%", 0.045))
	experiment.AddBandit(NewBandit("3%", 0.030))
	experiment.AddBandit(NewBandit("2%", 0.020))
	experiment.AddBandit(NewBandit("3.5%", 0.035))

	winner := pickOptimalVariant(experiment, 50000)
	fmt.Printf("%v\n", winner)
}
