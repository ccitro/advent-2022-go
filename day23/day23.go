package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var board [][]bool

type Pos struct {
	x, y int
}

type Bounds struct {
	topLeft, bottomRight Pos
}

const (
	North     = iota
	NorthEast = iota
	East      = iota
	SouthEast = iota
	South     = iota
	SouthWest = iota
	West      = iota
	NorthWest = iota
)

var directions = map[int]Pos{
	North:     {0, -1},
	NorthEast: {1, -1},
	East:      {1, 0},
	SouthEast: {1, 1},
	South:     {0, 1},
	SouthWest: {-1, 1},
	West:      {-1, 0},
	NorthWest: {-1, -1},
}

var movementScanning = map[int][]int{
	North: {NorthWest, North, NorthEast},
	South: {SouthWest, South, SouthEast},
	East:  {NorthEast, East, SouthEast},
	West:  {NorthWest, West, SouthWest},
}

var movementOrder = []int{North, South, West, East}

func loadPuzzle(file *os.File) {
	width := -1
	height := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		height++
		if width == -1 {
			width = len(line)
		}
	}

	sideBuffer := 200
	board = make([][]bool, height+sideBuffer*2)
	for i := range board {
		board[i] = make([]bool, width+sideBuffer*2)
	}

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x := 0; x < len(line); x++ {
			if line[x] == '#' {
				board[y+sideBuffer][x+sideBuffer] = true
			}
		}
		y++
	}
}

func findBounds() Bounds {
	minY := int(9e9)
	maxY := -1
	minX := int(9e9)
	maxX := -1
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] {
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
			}
		}
	}
	return Bounds{Pos{minX, minY}, Pos{maxX, maxY}}
}

func printBoard() {
	bounds := findBounds()
	for y := bounds.topLeft.y; y <= bounds.bottomRight.y; y++ {
		for x := bounds.topLeft.x; x <= bounds.bottomRight.x; x++ {
			if board[y][x] {
				print("#")
			} else {
				print(".")
			}
		}
		println()
	}
	println()
}

func getDirectionLabel(direction int) string {
	switch direction {
	case North:
		return "North"
	case NorthEast:
		return "NorthEast"
	case East:
		return "East"
	case SouthEast:
		return "SouthEast"
	case South:
		return "South"
	case SouthWest:
		return "SouthWest"
	case West:
		return "West"
	case NorthWest:
		return "NorthWest"
	}
	return "Unknown"
}

func determineDesiredPos(currentPos Pos) Pos {
	// fmt.Printf("Determining desired position for (%d, %d)\n", currentPos.x, currentPos.y)

	elfNearby := false
	for _, direction := range directions {
		scanPos := Pos{currentPos.x + direction.x, currentPos.y + direction.y}
		if board[scanPos.y][scanPos.x] {
			elfNearby = true
			break
		}
	}
	if !elfNearby {
		// fmt.Printf("No elf nearby, staying put\n")
		return currentPos
	}

	for _, direction := range movementOrder {
		// fmt.Printf("Checking %s\n", getDirectionLabel(direction))
		sideHasElf := false
		for _, scanDirection := range movementScanning[direction] {
			scanPos := Pos{currentPos.x + directions[scanDirection].x, currentPos.y + directions[scanDirection].y}
			if board[scanPos.y][scanPos.x] {
				// fmt.Printf("Found elf at (%d, %d)\n", scanPos.x, scanPos.y)
				sideHasElf = true
				break
			}
		}

		if !sideHasElf {
			// fmt.Printf("Found empty region to the %s\n", getDirectionLabel(direction))
			return Pos{currentPos.x + directions[direction].x, currentPos.y + directions[direction].y}
		}
	}

	// fmt.Printf("No empty region found, staying put\n")
	return currentPos
}

func moveElves() bool {
	destinationSquares := make(map[Pos][]Pos)
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if !board[y][x] {
				continue
			}

			currentPos := Pos{x, y}
			desiredPos := determineDesiredPos(currentPos)
			current, ok := destinationSquares[desiredPos]
			if !ok {
				destinationSquares[desiredPos] = []Pos{currentPos}
			} else {
				destinationSquares[desiredPos] = append(current, currentPos)
			}
		}
	}

	moved := false
	for destination, elves := range destinationSquares {
		if len(elves) > 1 {
			continue
		}

		if destination == elves[0] {
			continue
		}

		moved = true
		board[elves[0].y][elves[0].x] = false
		board[destination.y][destination.x] = true
	}

	movementOrder = append(movementOrder[1:], movementOrder[0])

	return moved
}

func part1() {
	println("== Initial State ==")
	printBoard()
	roundsRemaining := 10
	currentRound := 0

	for currentRound < roundsRemaining {
		currentRound++
		moveElves()

		fmt.Printf("== End of Round %d ==\n", currentRound)
		printBoard()
	}

	bounds := findBounds()
	emptyCount := 0
	for y := bounds.topLeft.y; y <= bounds.bottomRight.y; y++ {
		for x := bounds.topLeft.x; x <= bounds.bottomRight.x; x++ {
			if !board[y][x] {
				emptyCount++
			}
		}
	}
	println(emptyCount)

}

func part2() {
	currentRound := 0

	moved := true
	for moved {
		currentRound++
		moved = moveElves()
	}

	println(currentRound)
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
