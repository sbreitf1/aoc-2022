package main

// https://adventofcode.com/2022/day/7

import (
	"aoc/helper"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")
	commands := parseCommands(lines)

	root := &directory{Name: "/"}
	currentDir := root
	for _, cmd := range commands {
		switch cmd := cmd.(type) {
		case *cdCommand:
			if cmd.DirName == "/" {
				currentDir = root
			} else if cmd.DirName == ".." {
				currentDir = currentDir.Parent
				if currentDir == nil {
					helper.ExitWithMessage("cannot cd to parent of root")
				}
			} else {
				currentDir = currentDir.GetOrAddDir(cmd.DirName)
			}
		case *lsCommand:
			for _, dirName := range cmd.Dirs {
				currentDir.AddDir(dirName)
			}
			for _, file := range cmd.Files {
				currentDir.AddFile(file.Name, file.Size)
			}
		default:
			helper.ExitWithMessage("unexpected command %T", cmd)
		}
	}

	var sum1 int64
	for _, d := range root.GetDirsWithMaxTotalSize(100000) {
		sum1 += d.RecursiveSize
	}

	minDeleteSize := root.RecursiveSize - (70000000 - 30000000)
	deleteCandidates := root.GetDirsWithMinTotalSize(minDeleteSize)
	sort.Slice(deleteCandidates, func(i, j int) bool {
		return deleteCandidates[i].RecursiveSize < deleteCandidates[j].RecursiveSize
	})

	fmt.Println("-> part 1:", sum1)
	fmt.Println("-> part 2:", deleteCandidates[0].RecursiveSize)
}

type cdCommand struct {
	DirName string
}

type lsCommand struct {
	Dirs  []string
	Files []struct {
		Name string
		Size int64
	}
}

func parseCommands(lines []string) []interface{} {
	commands := make([]interface{}, 0)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}

		if strings.HasPrefix(l, "$ cd ") {
			commands = append(commands, &cdCommand{l[5:]})
			continue
		}
		if l == "$ ls" {
			commands = append(commands, &lsCommand{
				Dirs: make([]string, 0),
				Files: make([]struct {
					Name string
					Size int64
				}, 0),
			})
			continue
		}

		// probably ls output entry
		if len(commands) == 0 {
			helper.ExitWithMessage("command output %q found before first command", l)
		}
		lastCommand, ok := commands[len(commands)-1].(*lsCommand)
		if !ok {
			helper.ExitWithMessage("command output %q found after unexpected command", l)
		}

		parts := strings.SplitN(l, " ", 2)
		if len(parts) != 2 {
			helper.ExitWithMessage("unsupported command output %q", l)
		}
		if parts[0] == "dir" {
			lastCommand.Dirs = append(lastCommand.Dirs, parts[1])
			continue
		}
		fileSize, err := strconv.ParseInt(parts[0], 10, 64)
		helper.ExitOnError(err, "invalid file size of %q", parts[1])
		lastCommand.Files = append(lastCommand.Files, struct {
			Name string
			Size int64
		}{parts[1], fileSize})
	}
	return commands
}

type directory struct {
	Parent        *directory
	Name          string
	Dirs          []*directory
	Files         []*file
	RecursiveSize int64
}

func (d *directory) AddDir(name string) {
	for _, existingDir := range d.Dirs {
		if existingDir.Name == name {
			// directory already exists
			return
		}
	}
	d.Dirs = append(d.Dirs, &directory{Parent: d, Name: name})
}

func (d *directory) GetOrAddDir(name string) *directory {
	for _, existingDir := range d.Dirs {
		if existingDir.Name == name {
			return existingDir
		}
	}
	newDir := &directory{Parent: d, Name: name}
	d.Dirs = append(d.Dirs, newDir)
	return newDir
}

func (d *directory) AddFile(name string, size int64) {
	for _, existingFile := range d.Files {
		if existingFile.Name == name {
			if existingFile.Size != size {
				helper.ExitWithMessage("conflicting file sizes %d and %d found for file %s", size, existingFile.Size, name)
			}
			// file already exists
			return
		}
	}
	d.Files = append(d.Files, &file{Parent: d, Name: name, Size: size})
	for current := d; current != nil; current = current.Parent {
		current.RecursiveSize += size
	}
}

func (d *directory) GetDirsWithMaxTotalSize(maxSize int64) []*directory {
	dirs := make([]*directory, 0)
	if d.RecursiveSize <= maxSize {
		dirs = append(dirs, d)
	}
	for _, subDir := range d.Dirs {
		dirs = append(dirs, subDir.GetDirsWithMaxTotalSize(maxSize)...)
	}
	return dirs
}

func (d *directory) GetDirsWithMinTotalSize(minSize int64) []*directory {
	dirs := make([]*directory, 0)
	if d.RecursiveSize >= minSize {
		dirs = append(dirs, d)
	}
	for _, subDir := range d.Dirs {
		dirs = append(dirs, subDir.GetDirsWithMinTotalSize(minSize)...)
	}
	return dirs
}

type file struct {
	Parent *directory
	Name   string
	Size   int64
}
