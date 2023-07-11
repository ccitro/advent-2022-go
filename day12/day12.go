package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type point struct {
	x, y int
}

type path struct {
	points []point
}

type heightmap struct {
	terrain [][]int
	start   point
	end     point
}

var dirs = [4]point{
	{x: 0, y: 1},
	{x: 1, y: 0},
	{x: 0, y: -1},
	{x: -1, y: 0},
}

func loadHeightmap(file *os.File) *heightmap {
	hm := heightmap{}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, v := range line {
			if v == 'S' {
				hm.start = point{x: i, y: y}
				v = 'a'
			} else if v == 'E' {
				hm.end = point{x: i, y: y}
				v = 'z'
				// v = 'l'
			}

			row[i] = int(v - 'a')
		}
		hm.terrain = append(hm.terrain, row)
		y++
	}

	return &hm
}

func pathContainsPoint(p *path, pt point) bool {
	for _, v := range p.points {
		if v == pt {
			return true
		}
	}
	return false
}

func findShortestPath(hm *heightmap) *path {
	pathsLeft := []*path{{points: []point{hm.start}}}
	maxX := len(hm.terrain[0])
	maxY := len(hm.terrain)
	seenPoints := make(map[point]bool)

	for len(pathsLeft) > 0 {
		thisPath := pathsLeft[0]
		pathsLeft = append([]*path{}, pathsLeft[1:]...)

		currentPoint := thisPath.points[len(thisPath.points)-1]
		currentHeight := hm.terrain[currentPoint.y][currentPoint.x]
		if currentPoint == hm.end {
			return thisPath
		}
		seenPoints[currentPoint] = true

		for _, dir := range dirs {
			nextPoint := point{x: currentPoint.x + dir.x, y: currentPoint.y + dir.y}
			if nextPoint.x < 0 || nextPoint.x >= maxX || nextPoint.y < 0 || nextPoint.y >= maxY {
				continue
			}

			if seenPoints[nextPoint] {
				continue
			}

			if pathContainsPoint(thisPath, nextPoint) {
				continue
			}

			nextHeight := hm.terrain[nextPoint.y][nextPoint.x]
			if nextHeight-currentHeight > 1 {
				continue
			}

			newPoints := make([]point, len(thisPath.points)+1)
			copy(newPoints, thisPath.points)
			newPoints[len(newPoints)-1] = nextPoint
			pathsLeft = append(pathsLeft, &path{points: newPoints})
		}
	}

	panic("No path found")
}

func (hm *heightmap) print() {
	for y, row := range hm.terrain {
		for x, v := range row {
			if x == hm.start.x && y == hm.start.y {
				print("S")
			} else if x == hm.end.x && y == hm.end.y {
				print("E")
			} else {
				print(string(rune(v + 'a')))
			}
		}
		println()
	}
	println()
}

func part1(file *os.File) {
	hm := loadHeightmap(file)
	hm.print()

	path := findShortestPath(hm)
	fmt.Printf("Found path of length %d\n", len(path.points))
}

func part2(file *os.File) {
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
