package main

import (
	"code.google.com/p/gostat/stat"
	"fmt"
	"math"
	"sort"
)

const MeanIterations = 10000

type Bandit struct {
	Name  string
	alpha float64
	beta  float64

	Observations int64
	Rewards      int64
	Chosen       int64
}

func NewBandit(name string, successes int64, total int64) *Bandit {
	bandit := &Bandit{name, 0, 0, total, successes, 0}
	bandit.updateBetaParams()

	return bandit
}

func (b *Bandit) Observe() float64 {
	return stat.NextBeta(b.alpha, b.beta)
}

// The arithmetic mean and standard deviation, via sampling.
func (b *Bandit) Mean() (mean float64, stdDev float64, obs []float64) {
	obs = make([]float64, MeanIterations)
	var total float64
	for i := 0; i < MeanIterations; i++ {
		obs[i] = b.Observe()
		total += b.Observe()
	}

	mean = total / float64(MeanIterations)
	for i := 0; i < MeanIterations; i++ {
		stdDev += math.Pow(obs[i]-mean, 2)
	}

	stdDev = math.Sqrt(stdDev / float64(MeanIterations))

	return
}

func (b *Bandit) updateBetaParams() {
	b.alpha = 1.0 + float64(b.Rewards)
	b.beta = 1.0 + float64(b.Observations) - float64(b.Rewards)
}

func (b *Bandit) String() string {
	mean, stddev, _ := b.Mean()
	return fmt.Sprintf("%24s succ %10d trials %10d\tmean %3.6f +/- %3.6f",
		b.Name,
		b.Rewards,
		b.Observations,
		mean,
		stddev,
	)
}

// Relative and absolute difference in the
// two arms' conversion rate; the relative
// difference is relative to the other arm.
func (b *Bandit) Compare(other *Bandit) (relativeDifference []float64, absoluteDifference []float64) {
	// Sample each arm.
	_, _, ourSamples := b.Mean()
	_, _, theirSamples := other.Mean()

	relTmp := make([]float64, len(ourSamples))
	tmp := make([]float64, len(ourSamples))
	for i, sample := range(ourSamples) {
		relTmp[i] = (sample-theirSamples[i])/theirSamples[i]
		tmp[i] = sample - theirSamples[i]
	}

	sort.Float64s(tmp)
	sort.Float64s(relTmp)
	i5 := int(0.05 * float64(MeanIterations))
	i50 := int(0.5 * float64(MeanIterations))
	i95 := int(0.95 * float64(MeanIterations))

	relativeDifference = []float64{relTmp[i5], relTmp[i50], relTmp[i95]}
	absoluteDifference = []float64{tmp[i5], tmp[i50], tmp[i95]}
	return
}
