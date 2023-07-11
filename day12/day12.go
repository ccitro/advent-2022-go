package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type path = []int

type heightmap struct {
	terrain []int
	maxX    int
	start   int
	end     int
}

func coordsToIndex(x, y, maxX int) int {
	return y*maxX + x
}

func loadHeightmap(file *os.File) *heightmap {
	hm := heightmap{maxX: 0}
	y := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > hm.maxX {
			hm.maxX = len(line)
		}
		for i, v := range line {
			if v == 'S' {
				hm.start = coordsToIndex(i, y, hm.maxX)
				v = 'a'
			} else if v == 'E' {
				hm.end = coordsToIndex(i, y, hm.maxX)
				v = 'z'
				// v = 'l'
			}

			hm.terrain = append(hm.terrain, int(v-'a'))
		}
		y++
	}

	return &hm
}

func pathContainsPoint(p *path, pt int) bool {
	for _, v := range *p {
		if v == pt {
			return true
		}
	}
	return false
}

func findShortestPath(hm *heightmap) *path {
	defaultPath := path{hm.start}
	pathsLeft := []*path{&defaultPath}

	maxX := hm.maxX
	maxIndex := len(hm.terrain) - 1
	for len(pathsLeft) > 0 {
		fmt.Printf("Paths left: %d\n", len(pathsLeft))
		thisPath := pathsLeft[0]
		pathsLeft = pathsLeft[1:]

		currentPoint := (*thisPath)[len(*thisPath)-1]
		currentHeight := hm.terrain[currentPoint]
		if currentPoint == hm.end {
			return thisPath
		}

		for i := 0; i < 4; i++ {
			nextPoint := currentPoint
			if i == 0 {
				// east
				nextPoint++
				if nextPoint%maxX == 0 {
					continue
				}
			} else if i == 1 {
				// south
				nextPoint += maxX
				if nextPoint > maxIndex {
					continue
				}
			} else if i == 2 {
				// west
				nextPoint--
				if nextPoint%maxX == maxX-1 || nextPoint < 0 {
					continue
				}
			} else if i == 3 {
				// north
				nextPoint -= maxX
				if nextPoint < 0 {
					continue
				}
			}

			if pathContainsPoint(thisPath, nextPoint) {
				continue
			}

			nextHeight := hm.terrain[nextPoint]
			heightDifference := nextHeight - currentHeight
			if heightDifference > 1 {
				continue
			}

			newPath := make(path, len(*thisPath))
			copy(newPath, *thisPath)
			newPath = append(newPath, nextPoint)

			pathsLeft = append(pathsLeft, &newPath)
		}

	}

	panic("No path found")
}

func (hm *heightmap) print() {
	for i := 0; i < len(hm.terrain); i++ {
		if i == hm.start {
			print("S")
		} else if i == hm.end {
			print("E")
		} else {
			print(string(rune(hm.terrain[i] + 'a')))
		}

		if i%hm.maxX == hm.maxX-1 {
			println()
		}
	}
	println()
}

func part1(file *os.File) {
	hm := loadHeightmap(file)
	hm.print()

	path := findShortestPath(hm)
	fmt.Printf("Found path of length %d\n", len(*path))
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
