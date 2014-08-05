package main

import (
	"code.google.com/p/gostat/stat"
	"fmt"
	"math"
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
func (b *Bandit) Mean() (mean float64, stdDev float64) {
	obs := make([]float64, MeanIterations)
	var total float64
	for i := 0; i < MeanIterations; i++ {
		obs[i] = b.Observe()
		total += b.Observe()
	}

	mean = total/float64(MeanIterations);
	for i := 0; i < MeanIterations; i++ {
		stdDev += math.Pow(obs[i] - mean, 2)
	}

	stdDev = math.Sqrt(stdDev/float64(MeanIterations))

	return
}

func (b *Bandit) updateBetaParams() {
	b.alpha = 1.0 + float64(b.Rewards)
	b.beta = 1.0 + float64(b.Observations) - float64(b.Rewards)
}

func (b *Bandit) String() string {
	mean, stddev := b.Mean()
	return fmt.Sprintf("%24s succ %10d trials %10d mean %3.6f +/- %3.6f",
		b.Name,
		b.Rewards,
		b.Observations,
		mean,
		stddev,
	)
}
