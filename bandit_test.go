package main

import (
	"math/rand"
	"testing"
	"time"
)

const epsilon = 0.01

func InitRand() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestBetaBandit(t *testing.T) {
	InitRand()

	// Should simulate bernoulli with theta=0.5
	bandit := NewBandit("test", 50, 100)

	mean, _ := bandit.Mean()
	if !(mean - epsilon <= 0.5 && 0.5 <= mean + epsilon) {
		t.Errorf("Expected a mean close to 0.5 but was %f", mean)
	}
}
