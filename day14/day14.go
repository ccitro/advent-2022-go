package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const AIR = 0
const ROCK = 1
const SAND = 2

type rockPath = []point

type point struct {
	x int
	y int
}

type board struct {
	grid    [][]int
	offsetX int
	source  point
}

func parseLine(line string) rockPath {
	pointStrs := strings.Split(line, " -> ")
	rockPath := make(rockPath, len(pointStrs))
	for i, v := range pointStrs {
		coords := strings.Split(v, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		rockPath[i] = point{x: x, y: y}
	}

	return rockPath
}

func readBoard(file *os.File, hasFloor bool) board {
	sourceX := 500
	// first, extract the rock paths from the file
	scanner := bufio.NewScanner(file)
	rockPaths := make([]rockPath, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		rockPaths = append(rockPaths, parseLine(line))
	}

	// now, determine the size of the board
	// the intro example had an active range that didn't start until 400+,
	// so we calculate an offset to keep the grid size small
	minX := 9999999
	maxX := -1
	maxY := -1

	for _, v := range rockPaths {
		for _, p := range v {
			if p.x < minX {
				minX = p.x
			}
			if p.x > maxX {
				maxX = p.x
			}
			if p.y > maxY {
				maxY = p.y
			}
		}
	}

	if hasFloor {
		// for part 2, there is an "infinite" floor 2 below the lowest rock
		maxY += 2

		// sand can only effectively fall a horizontal distance equal to the vertical distance from the source
		furthestLeftSand := sourceX - maxY
		if furthestLeftSand < minX {
			minX = furthestLeftSand
		}
		furthestRightSand := sourceX + maxY
		if furthestRightSand > maxX {
			maxX = furthestRightSand
		}

		// add a path for the floor
		floorPath := make(rockPath, 2)
		floorPath[0] = point{x: minX, y: maxY}
		floorPath[1] = point{x: maxX, y: maxY}
		rockPaths = append(rockPaths, floorPath)
	}

	// now that we know the active area, we can create the grid
	width := maxX - minX + 1
	offsetX := minX
	grid := make([][]int, maxY+1)
	for i := 0; i < len(grid); i++ {
		grid[i] = make([]int, width)
	}

	// now, fill in the rocks, path by path
	for _, v := range rockPaths {
		for i := 1; i < len(v); i++ {
			// draw a line from the prior point to this point
			priorPoint := v[i-1]
			p := v[i]

			// this is a horizontal line
			if priorPoint.y == p.y {
				// determine if we're drawing from left to right or right to left
				minX := priorPoint.x
				maxX := p.x
				if priorPoint.x > p.x {
					minX = p.x
					maxX = priorPoint.x
				}
				for i := minX; i <= maxX; i++ {
					grid[p.y][i-offsetX] = ROCK
				}
			} else {
				// same as above, but for vertical lines
				minY := priorPoint.y
				maxY := p.y
				if priorPoint.y > p.y {
					minY = p.y
					maxY = priorPoint.y
				}
				for i := minY; i <= maxY; i++ {
					grid[i][p.x-offsetX] = ROCK
				}
			}
		}
	}

	source := point{x: sourceX - offsetX, y: 0}
	return board{grid: grid, offsetX: offsetX, source: source}
}

func (b board) print() {
	for x, v := range b.grid {
		for y, v2 := range v {
			if x == b.source.y && y == b.source.x {
				fmt.Printf("+")
				continue
			}
			switch v2 {
			case AIR:
				fmt.Printf(".")
			case ROCK:
				fmt.Printf("#")
			case SAND:
				fmt.Printf("o")
			}
		}
		fmt.Printf("\n")
	}

	println()
}

func addSand(b *board) bool {
	point := b.source
	// check if the source point is already sand
	if b.grid[point.y][point.x] == SAND {
		return false
	}

	for {
		if point.y >= len(b.grid)-1 {
			return false
		}

		// check if block below is air
		if b.grid[point.y+1][point.x] == AIR {
			point.y++
			continue
		}

		// next option is go go down and left
		// if we're at the left edge, then the sand will go to the void
		if point.x == 0 {
			return false
		}

		// check if down and left is air
		if b.grid[point.y+1][point.x-1] == AIR {
			point.x--
			point.y++
			continue
		}

		// next option is go go down and right
		// if we're at the right edge, then the sand will go to the void
		if point.x == len(b.grid[0])-1 {
			return false
		}

		// check if down and right is air
		if b.grid[point.y+1][point.x+1] == AIR {
			point.x++
			point.y++
			continue
		}

		// if we get here, then the sand settles here
		b.grid[point.y][point.x] = SAND
		return true
	}
}

func part1(file *os.File) {
	board := readBoard(file, false)

	sandCount := 0

	for {
		// board.print()
		addedSand := addSand(&board)
		if !addedSand {
			break
		}
		sandCount++
	}
	board.print()

	println(sandCount)
}

func part2(file *os.File) {
	board := readBoard(file, true)

	sandCount := 0

	for {
		// board.print()
		addedSand := addSand(&board)
		if !addedSand {
			break
		}
		sandCount++
	}
	board.print()

	println(sandCount)
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
	method(file)
}
