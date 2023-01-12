package main

// https://adventofcode.com/2022/day/16

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	nw := readValveNetwork(lines)

	maxPressureReleasePart1 := findMaxPressureReleasePart1(nw, "AA", "", 29, make(map[string]bool))

	fmt.Println("-> part 1:", maxPressureReleasePart1)
	//fmt.Println("-> part 2:", deleteCandidates[0].RecursiveSize)
}

type valveNetwork struct {
	valveRates         map[string]int
	valveJunctions     map[string][]string
	openableValveCount int
}

func readValveNetwork(lines []string) *valveNetwork {
	pattern := regexp.MustCompile(`^Valve ([A-Z]+) has flow rate=(\d+); tunnels? leads? to valves?\s+(.*)$`)
	nw := &valveNetwork{
		valveRates:     make(map[string]int),
		valveJunctions: make(map[string][]string),
	}
	for i, l := range lines {
		m := pattern.FindStringSubmatch(l)
		if len(m) != 4 {
			helper.ExitWithMessage("line %d %q did not match", i+1, l)
		}
		valveName := m[1]
		flowRate, _ := strconv.Atoi(m[2])
		parts := strings.Split(m[3], ",")
		nw.valveRates[valveName] = flowRate
		for _, p := range parts {
			nw.valveJunctions[valveName] = append(nw.valveJunctions[valveName], strings.TrimSpace(p))
		}
		if flowRate > 0 {
			nw.openableValveCount++
		}
	}
	return nw
}

func findMaxPressureReleasePart1(nw *valveNetwork, currentValve, previousValve string, remainingTime int, openValves map[string]bool) int {
	if remainingTime == 0 {
		// time is up
		return 0
	}
	if len(openValves) >= nw.openableValveCount {
		// no valves left to open
		return 0
	}

	maxPressureRelease := 0

	valveRate, ok := nw.valveRates[currentValve]
	if !ok {
		helper.ExitWithMessage("no valve rate known for %s", currentValve)
	}
	if valveRate > 0 && !openValves[currentValve] {
		// check solution with opening this valve:
		openValves[currentValve] = true
		// cut previous valve to allow going back after doing something here
		maxPressureRelease = valveRate*remainingTime + findMaxPressureReleasePart1(nw, currentValve, "", remainingTime-1, openValves)
		delete(openValves, currentValve)
	}

	valveJunctions, ok := nw.valveJunctions[currentValve]
	if !ok {
		helper.ExitWithMessage("no valve junctions known for %s", currentValve)
	}
	for _, nextValve := range valveJunctions {
		if nextValve == previousValve {
			// do not visit previous valve
			continue
		}
		// check for every junction
		pressureRelease := findMaxPressureReleasePart1(nw, nextValve, currentValve, remainingTime-1, openValves)
		if pressureRelease > maxPressureRelease {
			maxPressureRelease = pressureRelease
		}
	}
	return maxPressureRelease
}

//TODO remember visited nodes after last pressure release and deny visiting them again until next pressure release
//TODO compute best possible pressure release (all valves opened now) and compare to best known value
