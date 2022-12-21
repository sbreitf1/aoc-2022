package main

// https://adventofcode.com/2022/day/3

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	var prioritySum int
	for i, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		if len(l)%2 != 0 {
			helper.ExitWithMessage("line %d has invalid length", i+1)
		}

		commonRunes := findCommonRunes(l[:len(l)/2], l[len(l)/2:])
		if len(commonRunes) != 1 {
			helper.ExitWithMessage("line %d has %d runes in common", i+1)
		}
		prioritySum += getRunePriority(commonRunes[0])
	}

	var badgePrioritySum int
	for i := 0; i < len(lines)-2; i += 3 {
		commonRunes := findCommonRunes(lines[i], lines[i+1])
		commonRunes = findCommonRunes(string(commonRunes), lines[i+2])
		if len(commonRunes) != 1 {
			helper.ExitWithMessage("badge at line %d has %d runes in common", i+1)
		}
		badgePrioritySum += getRunePriority(commonRunes[0])
	}

	fmt.Println("-> part 1:", prioritySum)
	fmt.Println("-> part 2:", badgePrioritySum)
}

func findCommonRunes(str1, str2 string) []rune {
	m1 := make(map[rune]bool)
	for _, r := range str1 {
		m1[r] = true
	}
	m1and2 := make(map[rune]bool)
	for _, r := range str2 {
		if _, ok := m1[r]; ok {
			m1and2[r] = true
		}
	}
	commonRunes := make([]rune, 0)
	for r := range m1and2 {
		commonRunes = append(commonRunes, r)
	}
	return commonRunes
}

func getRunePriority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r) - int('a') + 1
	}
	if r >= 'A' && r <= 'Z' {
		return int(r) - int('A') + 27
	}
	panic(fmt.Sprintf("unsupported rune %q", r))
}
