package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Experiment struct {
	Name    string
	Bandits []*Bandit
}

type Result struct {
	Experiment              *Experiment
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

func (e *Experiment) pickOptimalVariant(iterations int) Result {
	rand.Seed(time.Now().UTC().UnixNano())

	allObservations := make([][]float64, iterations)
	bandits := len(e.Bandits)

	maxObs := 0.0
	maxVariantIndex := 0
	valueDist := make([]float64, iterations)
	counts := make([]int, bandits)

	for i := 0; i < iterations; i++ {
		// Sample every arm.
		observations := make([]float64, bandits)
		for j := 0; j < bandits; j++ {
			observations[j] = e.Bandits[j].Observe()
			if j == 0 || maxObs < observations[j] {
				maxVariantIndex = j
				maxObs = observations[j]
			}
		}

		allObservations[i] = observations
		counts[maxVariantIndex]++
	}

	maxCount := 0
	for i := 0; i < bandits; i++ {
		if i == 0 || maxCount < counts[i] {
			maxCount = counts[i]
			maxVariantIndex = i
		}
	}

	maxVariant := e.Bandits[maxVariantIndex]
	// Calculate the PVR remaining as the 95th percentile of the
	// posterior distribution of (t_max - t*)/t*, where t_max is the largest
	// observed arm sample for a given round of sampling, and t* is
	// the observation in the same round for the arm chosen as most likely to be optimal.
	for i := 0; i < iterations; i++ {
		observations := allObservations[i]

		// Find the maximal observation.
		for j := 0; j < bandits; j++ {
			if (j == 0) || (maxObs < observations[j]) {
				maxObs = observations[j]
			}
		}

		// Append (t_max-t*/t*) to our value distribution.
		optimalArmObs := observations[maxVariantIndex]
		valueDist[i] = (maxObs - optimalArmObs) / optimalArmObs
	}

	// Take the value at the 95th percentile.
	sort.Float64s(valueDist)
	p95th := (0.95 * float64(iterations))
	pvr := valueDist[int(p95th)]

	return Result{e, maxVariant, maxObs, pvr, iterations}
}

func (e *Experiment) String() string {
	return fmt.Sprintf("experiment \"%s\", %d arms", e.Name, len(e.Bandits))
}

func (r Result) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%s\n", r.Experiment.String())
	fmt.Fprintf(&buffer, "observations: %d\npotential value remaining: %f\n\n", r.Observations, r.PotentialValueRemaining)
	fmt.Fprintf(&buffer, "win:\t%s\n", r.Optimal.String())

	// TODO: be smart about finding control
	if r.Optimal != r.Experiment.Bandits[0] {
		relDiff, absDiff := r.Optimal.Compare(r.Experiment.Bandits[0])
		fmt.Fprintf(&buffer, "\n\t         ntile: %8d %8d %8d\n", 5, 50, 95)
		fmt.Fprintf(&buffer, "\trel to control: %3.6f %3.6f %3.6f\n", relDiff[0], relDiff[1], relDiff[2])
		fmt.Fprintf(&buffer, "\t           abs: %3.6f %3.6f %3.6f\n\n", absDiff[0], absDiff[1], absDiff[2])
	}

	for i := 0; i < len(r.Experiment.Bandits); i++ {
		bandit := r.Experiment.Bandits[i]
		if bandit != r.Optimal {
			fmt.Fprintf(&buffer, "arm:\t%s\n", bandit.String())
		}
	}
	return buffer.String()
}
