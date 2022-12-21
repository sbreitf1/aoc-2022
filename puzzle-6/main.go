package main

// https://adventofcode.com/2022/day/6

import (
	"aoc/helper"
	"fmt"
)

func main() {
	str := helper.ReadString("input.txt")

	var startOfPacketMarkerPos int
	for i := 4; i < len(str); i++ {
		if isMarker(str[i-4:i], 4) {
			startOfPacketMarkerPos = i
			break
		}
	}

	if startOfPacketMarkerPos < 4 {
		helper.ExitWithMessage("no start-of-packet-marker found")
	}

	var startOfMessageMarkerPos int
	for i := 14; i < len(str); i++ {
		if isMarker(str[i-14:i], 14) {
			startOfMessageMarkerPos = i
			break
		}
	}

	if startOfMessageMarkerPos < 14 {
		helper.ExitWithMessage("no start-of-message-marker found")
	}

	fmt.Println("-> part 1:", startOfPacketMarkerPos)
	fmt.Println("-> part 2:", startOfMessageMarkerPos)
}

func isMarker(str string, expectedLen int) bool {
	if len(str) != expectedLen {
		helper.ExitWithMessage("marker %q needs to have exactly %d characters", str, expectedLen)
	}
	usedRunes := make(map[rune]bool)
	for _, r := range str {
		usedRunes[r] = true
	}
	return len(usedRunes) == expectedLen
}
