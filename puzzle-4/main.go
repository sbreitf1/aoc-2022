package main

// https://adventofcode.com/2022/day/4

import (
	"aoc/helper"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	var fullOverlapCount, overlapCount int
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		range1, range2 := parseRanges(l)
		if areFullyOverlapping(range1, range2) {
			fullOverlapCount++
		}
		if areOverlapping(range1, range2) {
			overlapCount++
		}
	}

	fmt.Println("-> part 1:", fullOverlapCount)
	fmt.Println("-> part 2:", overlapCount)
}

type sectionRange struct {
	From, To int
}

func parseRanges(str string) (sectionRange, sectionRange) {
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		helper.ExitWithMessage("invalid entry %q", str)
	}
	parts1 := strings.Split(parts[0], "-")
	if len(parts1) != 2 {
		helper.ExitWithMessage("invalid entry %q", str)
	}
	from1, err := strconv.Atoi(parts1[0])
	helper.ExitOnError(err, "invalid entry %q", str)
	to1, err := strconv.Atoi(parts1[1])
	helper.ExitOnError(err, "invalid entry %q", str)
	parts2 := strings.Split(parts[1], "-")
	if len(parts1) != 2 {
		helper.ExitWithMessage("invalid entry %q", str)
	}
	from2, err := strconv.Atoi(parts2[0])
	helper.ExitOnError(err, "invalid entry %q", str)
	to2, err := strconv.Atoi(parts2[1])
	helper.ExitOnError(err, "invalid entry %q", str)
	return sectionRange{from1, to1}, sectionRange{from2, to2}
}

func areFullyOverlapping(range1, range2 sectionRange) bool {
	if range1.From >= range2.From && range1.To <= range2.To {
		return true
	}
	if range2.From >= range1.From && range2.To <= range1.To {
		return true
	}
	return false
}

func areOverlapping(range1, range2 sectionRange) bool {
	if range1.To < range2.From {
		return false
	}
	if range2.To < range1.From {
		return false
	}
	return true
}
