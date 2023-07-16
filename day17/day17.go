package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var jetPattern []int

type Shape = [][]int

var shapes = []Shape{
	// ####
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},

	// .#.
	// ###
	// .#.
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},

	// ..#
	// ..#
	// ###
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},

	// #
	// #
	// #
	// #
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},

	// ##
	// ##
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
}

var shapeHeights = []int{1, 3, 3, 4, 2}

type Chamber struct {
	rocks                   [][]bool
	highestSettledPoint     int
	fallingPieceSeq         int
	fallingPieceLowestPoint int
}

var chamber Chamber

const CHAMBER_WIDTH = 7
const ARRAY_CAPACITY = 10000
const ROCK_START_BOT_BUFFER = 3
const ROCK_START_LEFT_BUFFER = 2

func loadPuzzle(file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		l := len(line)
		jetPattern = make([]int, l)

		for i, v := range line {
			if v == '>' {
				jetPattern[i] = 1
			} else {
				jetPattern[i] = -1
			}
		}
	}
}

func makeChamber() {
	chamber = Chamber{
		highestSettledPoint:     0,
		rocks:                   make([][]bool, ARRAY_CAPACITY),
		fallingPieceSeq:         -1,
		fallingPieceLowestPoint: 0,
	}
	for i := range chamber.rocks {
		chamber.rocks[i] = make([]bool, CHAMBER_WIDTH)
	}
}

func isValidPosition(shape Shape, x int, y int) bool {
	for _, v := range shape {
		if y+v[1] < 0 {
			return false
		}
		if x+v[0] < 0 || x+v[0] >= CHAMBER_WIDTH {
			return false
		}
		if chamber.rocks[y+v[1]][x+v[0]] {
			return false
		}
	}
	return true
}

func placeShape(shape Shape, x int, y int) {
	for _, v := range shape {
		chamber.rocks[y+v[1]][x+v[0]] = true
	}
}

func (c *Chamber) print() {
	maxHeight := c.highestSettledPoint
	if c.fallingPieceSeq != -1 {
		fallingPieceMaxRow := c.fallingPieceLowestPoint + shapeHeights[c.fallingPieceSeq] - 1
		if fallingPieceMaxRow > maxHeight {
			maxHeight = fallingPieceMaxRow
		}
	}

	for y := maxHeight; y >= 0; y-- {
		print("|")
		for x := 0; x < len(c.rocks[y]); x++ {
			if c.rocks[y][x] {
				print("#")
			} else {
				print(".")
			}
		}
		println("|")
	}
	print("+")
	for i := 0; i < len(c.rocks[0]); i++ {
		print("-")
	}
	println("+")
	println()
}

type PieceJetCombo struct {
	piece int
	jet   int
}

type ComboState struct {
	seenCount           int
	pieceCount          int
	highestSettledPoint int
}

func doSimulation(totalRockCount int) {
	reportingInterval := 200
	if totalRockCount > 10000 {
		reportingInterval = 10000000
	}
	startTime := time.Now()

	makeChamber()

	rocksCompleted := 0
	shapeSeq := -1
	jetSeq := -1
	cycleHeightAdded := 0
	states := make(map[PieceJetCombo]*ComboState)

	for rocksCompleted < totalRockCount {
		if rocksCompleted%reportingInterval == 0 && rocksCompleted > 0 {
			elapsedSeconds := time.Since(startTime).Seconds()
			if rocksCompleted > 0 {
				remainingSeconds := (elapsedSeconds / float64(rocksCompleted)) * float64(totalRockCount-rocksCompleted)
				fmt.Printf("Completed %d rocks in %f seconds, %f remaining\n", rocksCompleted, elapsedSeconds, remainingSeconds)
			}
		}
		rocksCompleted++

		shapeSeq++
		shape := shapes[shapeSeq%len(shapes)]
		shapeX := 2
		shapeY := chamber.highestSettledPoint + ROCK_START_BOT_BUFFER

		for {
			jetSeq++
			jet := jetPattern[jetSeq%len(jetPattern)]
			newShapeX := shapeX + jet
			if isValidPosition(shape, newShapeX, shapeY) {
				shapeX = newShapeX
			}

			newShapeY := shapeY - 1
			if newShapeY >= 0 && isValidPosition(shape, shapeX, newShapeY) {
				shapeY = newShapeY
			} else {
				placeShape(shape, shapeX, shapeY)
				highestShapeY := shapeY + shapeHeights[shapeSeq%len(shapes)]
				if highestShapeY > chamber.highestSettledPoint {
					chamber.highestSettledPoint = highestShapeY
				}

				// i was going to convert the chamber to an array of ints, where each int is row with bitwise flags for each column
				// collision checks and writing to the chamber could be done with bitwise AND and OR operators, but
				// even so doing as many iterations as part 2 requires would probably take too long and require
				// cyclical arrays to avoid having a chamber that takes up too much memory

				// instead, this code finds a repeating pattern in the order that pieces and jets come up,
				// and uses that to skip ahead in the simulation
				if cycleHeightAdded == 0 {
					combo := PieceJetCombo{piece: shapeSeq % len(shapes), jet: jetSeq % len(jetPattern)}
					state, ok := states[combo]
					if ok && state.seenCount == 2 {
						topIncrease := chamber.highestSettledPoint - state.highestSettledPoint
						pieceIncrease := rocksCompleted - state.pieceCount
						maxRepeatPossible := (totalRockCount - rocksCompleted) / pieceIncrease
						cycleHeightAdded = topIncrease * maxRepeatPossible
						rocksCompleted += pieceIncrease * maxRepeatPossible
						fmt.Printf("Skipping %d cycles by adding %d rocks which will increase height by %d\n", maxRepeatPossible, pieceIncrease*maxRepeatPossible, topIncrease*maxRepeatPossible)
					}

					if !ok {
						state = &ComboState{seenCount: 1, pieceCount: rocksCompleted, highestSettledPoint: chamber.highestSettledPoint}
						states[combo] = state
					} else {
						state.seenCount++
						state.pieceCount = rocksCompleted
						state.highestSettledPoint = chamber.highestSettledPoint
					}
				}
				break
			}

		}
	}

	println(chamber.highestSettledPoint + cycleHeightAdded)
	fmt.Printf("Completed in %f milliseconds\n", time.Since(startTime).Seconds()*1000)
}

func part1() {
	doSimulation(2022)
}

func part2() {
	doSimulation(1000000000000)
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
