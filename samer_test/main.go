package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/samepic"
)

const DefaultSampleCount = 100

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "samer samples [count]")
		fmt.Fprintln(os.Stderr, "Available samers:")
		fmt.Fprintln(os.Stderr, " - avghash")
		fmt.Fprintln(os.Stderr, " - colorprof")
		fmt.Fprintln(os.Stderr, " - squashcomp")
		fmt.Fprintln(os.Stderr, " - neuralnet[PATH] ([PATH] is a filepath)")
		os.Exit(1)
	}

	count := DefaultSampleCount
	if len(os.Args) == 4 {
		var err error
		count, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid count:", os.Args[3])
			os.Exit(1)
		}
	}

	var samer samepic.Samer
	switch os.Args[1] {
	case "avghash":
		samer = &samepic.AverageHash{}
	case "colorprof":
		samer = &samepic.ColorProf{}
	case "squashcomp":
		samer = &samepic.SquashComp{}
	default:
		if strings.HasPrefix(os.Args[1], "neuralnet") {
			path := os.Args[1][len("neuralnet"):]
			var err error
			samer, err = samepic.LoadNeuralSamer(path)
			if err != nil {
				essentials.Die(err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Unknown samer:", os.Args[1])
			os.Exit(1)
		}
	}
	samples, err := samepic.NewDirSamples(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read sample dir:", err)
		os.Exit(1)
	}
	fmt.Println("Rating...")
	pos, neg, err := samepic.Rate(samer, samples, samepic.DefaultManipulator, count)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to rate:", err)
		os.Exit(1)
	}
	fmt.Println("Positive rating:", pos)
	fmt.Println("Negative rating:", neg)
}
