package main

import (
	"fmt"
)

type Experiment struct {
	Name    string
	Bandits []*Bandit
}

type Result struct {
	Optimal                 *Bandit
	ExpectedValue           float64
	PotentialValueRemaining float64
	Observations            int
}

func NewExperiment(name string) *Experiment {
	return &Experiment{name, []*Bandit{}}
}

func (e *Experiment) AddBandit(arm *Bandit) {
	e.Bandits = append(e.Bandits, arm)
}

func (e *Experiment) optimalVariant() (maxVariant *Bandit, maxVariantSample float64) {
	var sample float64
	maxVariant = nil
	maxVariantSample = 0.0

	for _, variant := range e.Bandits {
		// Find the variant that generates the highest value from Thompson sampling.
		if maxVariant == nil {
			maxVariant = variant
			maxVariantSample = variant.estimatedConversionRate()
		} else {
			sample = variant.estimatedConversionRate()
			if sample > maxVariantSample {
				maxVariant = variant
				maxVariantSample = sample
			}
		}
	}

	maxVariant.Chosen++
	return
}

func pickOptimalVariant(experiment *Experiment, iterations int) Result {
	for i := 0; i < iterations; i++ {
		// Find the arm with highest expected reward for the next pull.
		maxVariant, _ := experiment.optimalVariant()

		// Pull the chosen arm, and update its estimated conversion rate based on the observation.
		maxVariant.observe()
		maxVariant.updateBetaParams()
	}

	maxVariant, maxVariantSample := experiment.optimalVariant()
	return Result{maxVariant, maxVariantSample, 0.0, iterations}
}

func (e *Experiment) String() string {
	return fmt.Sprintf("experiment %s, %d arms", e.Name, len(e.Bandits))
}

func (r Result) String() string {
	return fmt.Sprintf("winner with expected conversion %f after %d observations: %v", r.ExpectedValue, r.Observations, r.Optimal)
}
