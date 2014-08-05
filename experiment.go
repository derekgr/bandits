package main

import (
	"fmt"
	"sort"
	"bytes"
)

type Experiment struct {
	Name    string
	Bandits []*Bandit

	observations [][]float64
}

type Result struct {
	Experiment              *Experiment
	Optimal                 *Bandit
	ExpectedValue           float64
	PotentialValueRemaining float64
	Observations            int
}

func NewExperiment(name string) *Experiment {
	return &Experiment{name, []*Bandit{}, nil}
}

func (e *Experiment) AddBandit(arm *Bandit) {
	e.Bandits = append(e.Bandits, arm)
}

func (e *Experiment) optimalVariant() (maxVariant *Bandit, maxVariantIndex int, maxVariantSample float64, observations []float64) {
	var sample float64
	maxVariant = nil
	maxVariantSample = 0.0
	observations = make([]float64, len(e.Bandits))

	for i, variant := range e.Bandits {
		sample = variant.Observe()
		observations[i] = sample

		// Find the variant that generates the highest value from Thompson sampling.
		if (maxVariant == nil) || (sample > maxVariantSample) {
			maxVariant = variant
			maxVariantSample = sample
			maxVariantIndex = i
		}
	}

	maxVariant.Chosen++
	return
}

func pickOptimalVariant(experiment *Experiment, iterations int) Result {
	experiment.observations = make([][]float64, iterations)
	bandits := len(experiment.Bandits)

	maxObs := 0.0
	maxVariantIndex := 0
	valueDist := make([]float64, iterations)
	counts := make([]int, bandits)

	for i := 0; i < iterations; i++ {
		// Sample every arm.
		observations := make([]float64, bandits)
		for j := 0; j < bandits; j++ {
			observations[j] = experiment.Bandits[j].Observe()
			if (j == 0 || maxObs < observations[j]) {
				maxVariantIndex = j
				maxObs = observations[j]
			}
		}

		experiment.observations[i] = observations
		counts[maxVariantIndex]++
	}

	maxCount := 0
	for i := 0; i < bandits; i++ {
		if (i == 0 || maxCount < counts[i]) {
			maxCount = counts[i]
			maxVariantIndex = i
		}
	}

	maxVariant := experiment.Bandits[maxVariantIndex]
    // Calculate the PVR remaining as the 95th percentile of the
	// posterior distribution of (t_max - t*)/t*, where t_max is the largest
	// observed arm sample for a given round of sampling, and t* is
	// the observation in the same round for the arm chosen as most likely to be optimal.
	for i := 0; i < iterations; i++ {
		observations := experiment.observations[i]

		// Find the maximal observation.
		for j := 0; j < bandits; j++ {
			if (j == 0) || (maxObs < observations[j]) {
				maxObs = observations[j]
			}
		}

		// Append (t_max-t*/t*) to our value distribution.
		optimalArmObs := observations[maxVariantIndex]
		valueDist[i] = (maxObs - optimalArmObs)/optimalArmObs
	}

	// Take the value at the 95th percentile.
	sort.Float64s(valueDist)
	p95th := (0.95 * float64(iterations))
	pvr := valueDist[int(p95th)]

	return Result{experiment, maxVariant, maxObs, pvr, iterations}
}

func (e *Experiment) String() string {
	return fmt.Sprintf("experiment %s, %d arms", e.Name, len(e.Bandits))
}

func (r Result) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "observations: %d  pvr: %f\n", r.Observations, r.PotentialValueRemaining)
	fmt.Fprintf(&buffer, "win:\tname %s\tsucc %d\tobs %d\test conversion rate %f\n", r.Optimal.Name, r.Optimal.Rewards, r.Optimal.Observations, r.ExpectedValue)
	for i := 0; i < len(r.Experiment.Bandits); i++ {
		bandit := r.Experiment.Bandits[i]
		if bandit != r.Optimal {
			fmt.Fprintf(&buffer, "arm:\tname %s\tsucc %d\tobs %d\test conversion rate %f\n", bandit.Name, bandit.Rewards, bandit.Observations, bandit.Observe())
		}
	}
	return buffer.String()
}
