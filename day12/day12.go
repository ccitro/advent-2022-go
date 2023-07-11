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
	endX    int
	endY    int
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
				hm.endX = i
				hm.endY = y
				v = 'z'
			}

			hm.terrain = append(hm.terrain, int(v-'a'))
		}
		y++
	}

	return &hm
}

func reconstructPath(cameFrom map[int]int, current int, origin int) *path {
	maxLength := len(cameFrom)
	totalPath := path{current}
	for {
		if _, ok := cameFrom[current]; !ok {
			break
		}
		if current == origin {
			break
		}
		current = cameFrom[current]
		totalPath = append(totalPath, current)
		if len(totalPath) > maxLength {
			panic("Path too long")
		}
	}
	return &totalPath
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func heuristicCostEstimate(p int, hm *heightmap) int {
	x := p % hm.maxX
	y := p / hm.maxX

	return abs(x-hm.endX) + abs(y-hm.endY)
}

func getNeighbors(current int, hm *heightmap) []int {
	neighbors := []int{}
	maxX := hm.maxX
	maxIndex := len(hm.terrain) - 1
	currentHeight := hm.terrain[current]
	for i := 0; i < 4; i++ {
		nextPoint := current
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
		nextPointHeight := hm.terrain[nextPoint]
		heightDifference := nextPointHeight - currentHeight
		if heightDifference > 1 {
			continue
		}
		neighbors = append(neighbors, nextPoint)
	}

	return neighbors
}

func getNeighborWeight(current int, neighbor int, hm *heightmap) int {
	currentHeight := hm.terrain[current]
	neighborHeight := hm.terrain[neighbor]
	return 2 - (neighborHeight - currentHeight)
}

func A_Star(start int, goal int, hm *heightmap) *path {
	inspections := 0
	openSet := make(map[int]bool)
	openSet[start] = true

	cameFrom := make(map[int]int)
	gScore := make(map[int]int)
	gScore[start] = 0

	fScore := make(map[int]int)
	fScore[start] = heuristicCostEstimate(start, hm)

	for len(openSet) > 0 {
		current := -1
		for k, _ := range openSet {
			if current == -1 || fScore[k] < fScore[current] {
				current = k
			}
		}

		if current == goal {
			fmt.Printf("Inspected %d nodes\n", inspections)
			return reconstructPath(cameFrom, current, start)
		}

		delete(openSet, current)
		for _, neighbor := range getNeighbors(current, hm) {
			inspections++
			tentative_gScore := gScore[current] + getNeighborWeight(current, neighbor, hm)
			neighborScore := gScore[neighbor]
			if neighborScore == 0 {
				neighborScore = 999999999
			}

			if tentative_gScore < neighborScore {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentative_gScore
				fScore[neighbor] = gScore[neighbor] + heuristicCostEstimate(neighbor, hm)
				if !openSet[neighbor] {
					openSet[neighbor] = true
				}
			}
		}
	}

	panic("No path found")
}

func findShortestPath(hm *heightmap) *path {
	return A_Star(hm.start, hm.end, hm)
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
	fmt.Printf("Found path with %d steps\n", len(*path)-1)
	println()
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
