package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	experiment := NewExperiment("console")
	experiment.AddBandit(NewBandit("control-avatar", 42060, 496653))
	experiment.AddBandit(NewBandit("test-avatar", 15280, 166337))

	winner := pickOptimalVariant(experiment, 1000)
	fmt.Printf("%v\n", winner)
}
