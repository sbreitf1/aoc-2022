package main

// https://adventofcode.com/2022/day/16

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	lines := helper.ReadLines("input.txt")

	nw := readValveNetwork(lines)

	start1 := time.Now()
	maxPressureReleasePart1 := findMaxPressureReleasePart1(nw, &dude{currentValve: "AA", visitedValvesSinceLastOpen: map[string]bool{"AA": true}}, 29, make(map[string]bool))
	fmt.Println("took", time.Since(start1))

	fmt.Println("-> part 1:", maxPressureReleasePart1)
	//fmt.Println("-> part 2:", deleteCandidates[0].RecursiveSize)
}

type valveNetwork struct {
	valveRates         map[string]int
	valveJunctions     map[string][]string
	openableValveCount int
}

type dude struct {
	currentValve               string
	visitedValvesSinceLastOpen map[string]bool
}

func (d *dude) Enter(valve string) string {
	previousValve := d.currentValve
	d.currentValve = valve
	d.visitedValvesSinceLastOpen[valve] = true
	return previousValve
}

func (d *dude) LeaveToPrevious(previousValve string) {
	delete(d.visitedValvesSinceLastOpen, d.currentValve)
	d.currentValve = previousValve
}

func (d *dude) OpenValve() *dude {
	return &dude{
		currentValve:               d.currentValve,
		visitedValvesSinceLastOpen: map[string]bool{d.currentValve: true},
	}
}

func (d *dude) CanEnter(valve string) bool {
	return !d.visitedValvesSinceLastOpen[valve]
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

func findMaxPressureReleasePart1(nw *valveNetwork, dude1 *dude, remainingTime int, openValves map[string]bool) int {
	if remainingTime == 0 {
		// time is up
		return 0
	}
	if len(openValves) >= nw.openableValveCount {
		// no valves left to open
		return 0
	}

	maxPressureRelease := 0

	valveRate, ok := nw.valveRates[dude1.currentValve]
	if !ok {
		helper.ExitWithMessage("no valve rate known for %s", dude1.currentValve)
	}
	if valveRate > 0 && !openValves[dude1.currentValve] {
		// check solution with opening this valve:
		openValves[dude1.currentValve] = true
		// cut previous valve to allow going back after doing something here
		maxPressureRelease = valveRate*remainingTime + findMaxPressureReleasePart1(nw, dude1.OpenValve(), remainingTime-1, openValves)
		delete(openValves, dude1.currentValve)
	}

	valveJunctions, ok := nw.valveJunctions[dude1.currentValve]
	if !ok {
		helper.ExitWithMessage("no valve junctions known for %s", dude1.currentValve)
	}
	for _, nextValve := range valveJunctions {
		if !dude1.CanEnter(nextValve) {
			// do not visit previous valve
			continue
		}
		// check for every junction
		previousValve := dude1.Enter(nextValve)
		pressureRelease := findMaxPressureReleasePart1(nw, dude1, remainingTime-1, openValves)
		dude1.LeaveToPrevious(previousValve)
		if pressureRelease > maxPressureRelease {
			maxPressureRelease = pressureRelease
		}
	}
	return maxPressureRelease
}

//TODO remember visited nodes after last pressure release and deny visiting them again until next pressure release
//TODO compute best possible pressure release (all valves opened now) and compare to best known value
