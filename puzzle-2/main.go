package main

// https://adventofcode.com/2022/day/2

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	selectionPoints := map[string]int{
		"X": 1,
		"Y": 2,
		"Z": 3,
	}
	outcome1Points := map[string]map[string]int{
		"X": {
			"A": 3,
			"B": 0,
			"C": 6,
		},
		"Y": {
			"A": 6,
			"B": 3,
			"C": 0,
		},
		"Z": {
			"A": 0,
			"B": 6,
			"C": 3,
		},
	}
	outcome2Points := map[string]map[string]int{
		"A": { // rock
			"X": 0 + 3, // mine: scissors
			"Y": 3 + 1, // mine: rock
			"Z": 6 + 2, // mine: paper
		},
		"B": { // paper
			"X": 0 + 1, // mine: rock
			"Y": 3 + 2, // mine: paper
			"Z": 6 + 3, // mine: scissors
		},
		"C": { // scissors
			"X": 0 + 2, // mine: paper
			"Y": 3 + 3, // mine: scissors
			"Z": 6 + 1, // mine: rock
		},
	}

	var points1, points2 int
	for i, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		parts := strings.Split(l, " ")
		if len(parts) != 2 {
			helper.ExitWithMessage("line %d seems broken", (i + 1))
		}
		other := parts[0]
		mine := parts[1]
		pSel, ok := selectionPoints[mine]
		if !ok {
			helper.ExitWithMessage("unknown selection %q", mine)
		}
		mOutcome1, ok := outcome1Points[mine]
		if !ok {
			helper.ExitWithMessage("unknown selection %q", mine)
		}
		pOutcome1, ok := mOutcome1[other]
		if !ok {
			helper.ExitWithMessage("unknown other selection %q", other)
		}
		mOutcome2, ok := outcome2Points[other]
		if !ok {
			helper.ExitWithMessage("unknown other selection %q", other)
		}
		pOutcome2, ok := mOutcome2[mine]
		if !ok {
			helper.ExitWithMessage("unknown selection %q", mine)
		}
		points1 += pSel + pOutcome1
		points2 += pOutcome2
	}

	fmt.Println("-> part 1:", points1)
	fmt.Println("-> part 2:", points2)
}
