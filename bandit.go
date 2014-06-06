package main

import (
	"code.google.com/p/gostat/stat"
	"fmt"
)

type Bandit struct {
	Name  string
	alpha float64
	beta  float64

	Observations int64
	Rewards      int64
	Chosen       int64
	sampling     func() int64
}

func NewBandit(name string, conversionRate float64) *Bandit {
	return &Bandit{name, 1.0, 1.0, 0, 0, 0, stat.Bernoulli(conversionRate)}
}

func (b *Bandit) estimatedConversionRate() float64 {
	return stat.NextBeta(b.alpha, b.beta)
}

func (b *Bandit) updateBetaParams() {
	b.alpha = 1.0 + float64(b.Rewards)
	b.beta = 1.0 + float64(b.Observations) - float64(b.Rewards)
}

func (b *Bandit) observe() int64 {
	reward := b.sampling()

	b.Observations++
	b.Rewards += reward

	return reward
}

func (b *Bandit) String() string {
	return fmt.Sprintf("arm %s, chosen as optimal = %d, obs = %d, succ = %d, fail = %d",
		b.Name,
		b.Chosen,
		b.Observations,
		b.Rewards,
		b.Observations-b.Rewards)
}
