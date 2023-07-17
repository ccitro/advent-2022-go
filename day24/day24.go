package main

import (
	"bufio"
	"os"
	"strings"
)

var maze [][]string

type Pos struct {
	x int
	y int
}

var dirs = []Pos{
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
}

func loadPuzzle(file *os.File) {
	maze = make([][]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "####") || strings.HasSuffix(line, "####") {
			continue
		}

		cells := strings.Split(line, "")
		maze = append(maze, cells[1:len(cells)-1])
	}
}

func mod(x, m int) int {
	if x == 0 {
		return 0
	}

	if x > 0 {
		return x % m
	}

	for x < 0 {
		x += m
	}

	return x
}

func search(start Pos, exit Pos) int {
	step := 1
	height := len(maze)
	width := len(maze[0])

	positions := []Pos{start}
	for {
		nextPositions := make([]Pos, 0)
		for _, s := range positions {
			for _, d := range dirs {
				newPos := Pos{s.x + d.x, s.y + d.y}
				if newPos.x == exit.x && newPos.y == exit.y {
					return step
				}

				if newPos.x < 0 || newPos.x >= width || newPos.y < 0 || newPos.y >= height {
					continue
				}

				if maze[newPos.y][mod(newPos.x-step, width)] == ">" {
					continue
				}

				if maze[newPos.y][mod(newPos.x+step, width)] == "<" {
					continue
				}

				if maze[mod(newPos.y+step, height)][newPos.x] == "V" {
					continue
				}

				if maze[mod(newPos.y-step, height)][newPos.x] == "^" {
					continue
				}

				nextPositions = append(nextPositions, newPos)
			}
		}

		positions = nextPositions
		if len(positions) == 0 {
			positions = []Pos{start}
		}
		step++
	}

}

func part1() {
	height := len(maze)
	width := len(maze[0])
	start := Pos{-1, 0}
	exit := Pos{width - 1, height}

	steps := search(start, exit)
	println(steps)
}

func part2() {
}

func main() {
	filename := "input.txt"
	method := part1
	for _, v := range os.Args {
		if v == "part2" || v == "2" {
			method = part2
		}
		if strings.HasSuffix(v, ".txt") {
			filename = v
		}
	}

	file, _ := os.Open(filename)
	defer file.Close()
	loadPuzzle(file)
	method()
}
