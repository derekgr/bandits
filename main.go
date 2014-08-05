package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		experimentName = flag.String("name", "console", "Experiment name")
		iterations     = flag.Int("iter", 10000, "Iterations of sampling")
	)

	flag.Usage = func() {
		fmt.Printf("%s [--iter iterations] [--name experiment-name] [csv-file]\n\n"+
			"If csv-file is omitted, data will be read from STDIN.\n"+
			"The expected format is CSV, with one row per bandit, and each row being:\n"+
			"arm_name,arm_successes,arm_total_trials\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	experiment := NewExperiment(*experimentName)
	file := os.Stdin
	if len(flag.Args()) > 0 {
		var err error
		file, err = os.Open(flag.Args()[0])
		if err != nil {
			panic(err)
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			panic(fmt.Errorf("Malformed input line: %s", scanner.Text()))
		}
		successes, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		trials, err := strconv.Atoi(parts[2])
		if err != nil {
			panic(err)
		}

		experiment.AddBandit(NewBandit(parts[0], int64(successes), int64(trials)))
	}

	winner := experiment.pickOptimalVariant(*iterations)
	fmt.Printf("%s\n", winner.String())
}
