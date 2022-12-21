package main

// https://adventofcode.com/2022/day/5

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	patternMoveCommand = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)
)

func main() {
	lines := helper.ReadLines("input.txt")

	headerLines := make([]string, 0)
	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		headerLines = append(headerLines, l)
	}

	ship1 := newCargoShip(headerLines)
	ship2 := newCargoShip(headerLines)

	lines = lines[len(headerLines)+1:]
	for _, l := range lines {
		cmd := parseMoveCommand(l)
		for i := 0; i < cmd.Count; i++ {
			ship1.Move(cmd.From, cmd.To)
		}
		ship2.MoveMany(cmd.From, cmd.To, cmd.Count)
	}

	fmt.Println("-> part 1:", ship1.CurrentTopMessage())
	fmt.Println("-> part 2:", ship2.CurrentTopMessage())
}

type cargoShip struct {
	Stacks [][]rune
}

func newCargoShip(headerLines []string) *cargoShip {
	stackCount := getStackCount(headerLines[len(headerLines)-1])
	stacks := make([][]rune, stackCount)
	for i := 0; i < len(headerLines)-1; i++ {
		index := len(headerLines) - i - 2
		cargoLine := parseCargoLine(headerLines[index])
		if len(cargoLine) != stackCount {
			helper.ExitWithMessage("found cargo line with %d entries, but expected %d", len(cargoLine), stackCount)
		}
		for j, r := range cargoLine {
			if r != ' ' {
				if len(stacks[j]) < i {
					helper.ExitWithMessage("stack %d had %d items, but exptected %d", j+1, len(stacks[j]), i)
				}
				stacks[j] = append(stacks[j], r)
			}
		}
	}

	return &cargoShip{
		Stacks: stacks,
	}
}

func getStackCount(str string) int {
	var stackCount int
	parts := strings.Split(str, " ")
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}

		num, err := strconv.Atoi(p)
		helper.ExitOnError(err, "invalid stack header %q", str)
		if num != stackCount+1 {
			helper.ExitWithMessage("invalid stack header %q", str)
		}
		stackCount = num
	}
	return stackCount
}

func parseCargoLine(str string) []rune {
	cargoLine := make([]rune, 0)
	for i := 0; i < len(str)-2; i += 4 {
		if str[i] == '[' && str[i+2] == ']' {
			cargoLine = append(cargoLine, rune(str[i+1]))
		} else if str[i] == ' ' && str[i+1] == ' ' && str[i+2] == ' ' {
			cargoLine = append(cargoLine, ' ')
		} else {
			helper.ExitWithMessage("invalid cargo line %q", str)
		}
	}
	return cargoLine
}

func (ship *cargoShip) String() string {
	var highestStack int
	for _, s := range ship.Stacks {
		if len(s) > highestStack {
			highestStack = len(s)
		}
	}

	var sb strings.Builder
	for i := highestStack - 1; i >= 0; i-- {
		for j := 0; j < len(ship.Stacks); j++ {
			if j > 0 {
				sb.WriteRune(' ')
			}
			if len(ship.Stacks[j]) > i {
				sb.WriteRune('[')
				sb.WriteRune(ship.Stacks[j][i])
				sb.WriteRune(']')
			} else {
				sb.WriteString("   ")
			}
		}
		sb.WriteRune('\n')
	}
	for j := 0; j < len(ship.Stacks); j++ {
		if j > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteRune(' ')
		sb.WriteString(strconv.Itoa(j + 1))
		sb.WriteRune(' ')
	}
	return sb.String()
}

func (ship *cargoShip) CurrentTopMessage() string {
	var sb strings.Builder
	for j := 0; j < len(ship.Stacks); j++ {
		if len(ship.Stacks[j]) > 0 {
			sb.WriteRune(ship.Stacks[j][len(ship.Stacks[j])-1])
		} else {
			sb.WriteRune(' ')
		}
	}
	return sb.String()
}

func (ship *cargoShip) Move(from, to int) {
	if len(ship.Stacks[from-1]) == 0 {
		helper.ExitWithMessage("cannot move from empty stack %d", from)
	}
	ship.Stacks[to-1] = append(ship.Stacks[to-1], ship.Stacks[from-1][len(ship.Stacks[from-1])-1])
	ship.Stacks[from-1] = ship.Stacks[from-1][:len(ship.Stacks[from-1])-1]
}

func (ship *cargoShip) MoveMany(from, to int, count int) {
	if len(ship.Stacks[from-1]) < count {
		helper.ExitWithMessage("cannot move %d items from stack of size %d", count, from)
	}
	ship.Stacks[to-1] = append(ship.Stacks[to-1], ship.Stacks[from-1][len(ship.Stacks[from-1])-count:len(ship.Stacks[from-1])]...)
	ship.Stacks[from-1] = ship.Stacks[from-1][:len(ship.Stacks[from-1])-count]
}

type moveCommand struct {
	Count    int
	From, To int
}

func parseMoveCommand(str string) moveCommand {
	m := patternMoveCommand.FindStringSubmatch(str)
	if len(m) != 4 {
		helper.ExitWithMessage("invalid command %q", str)
	}
	count, _ := strconv.Atoi(m[1])
	from, _ := strconv.Atoi(m[2])
	to, _ := strconv.Atoi(m[3])
	return moveCommand{count, from, to}
}
