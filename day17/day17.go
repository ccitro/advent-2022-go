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
	rocks                   [][]int
	highestSettledPoint     int
	fallingPieceSeq         int
	fallingPieceLowestPoint int
	heightBelowFLoor        int
}

var chamber Chamber

const CHAMBER_WIDTH = 7
const LONGEST_PERFORATED_STRETCH = 1000
const ROCK_START_BOT_BUFFER = 3
const ROCK_START_LEFT_BUFFER = 2
const CONTENTS_SETTLED_ROCK = 1
const CONTENTS_FALLING_ROCK = 2

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
		rocks:                   make([][]int, LONGEST_PERFORATED_STRETCH),
		fallingPieceSeq:         -1,
		fallingPieceLowestPoint: 0,
		heightBelowFLoor:        0,
	}
	for i := range chamber.rocks {
		chamber.rocks[i] = make([]int, CHAMBER_WIDTH)
	}
}

func placeShape(shapeSeq int) {
	bottom := chamber.highestSettledPoint + ROCK_START_BOT_BUFFER
	left := ROCK_START_LEFT_BUFFER
	shape := shapes[shapeSeq]

	for _, v := range shape {
		x := v[0] + left
		y := v[1] + bottom
		chamber.rocks[y][x] = CONTENTS_FALLING_ROCK
	}

	chamber.fallingPieceSeq = shapeSeq
	chamber.fallingPieceLowestPoint = bottom
}

func applyJet(jetSeq int) bool {
	if chamber.fallingPieceSeq == -1 {
		return false
	}

	dir := jetPattern[jetSeq]
	canMoveHorizontally := true
	for y := chamber.fallingPieceLowestPoint; y < chamber.fallingPieceLowestPoint+shapeHeights[chamber.fallingPieceSeq]; y++ {
		for x, contents := range chamber.rocks[y] {
			if contents == CONTENTS_FALLING_ROCK {
				if (dir == -1 && x == 0) || (dir == 1 && x == len(chamber.rocks[y])-1) {
					canMoveHorizontally = false
					break
				}

				neighborContents := chamber.rocks[y][x+dir]
				if neighborContents == CONTENTS_SETTLED_ROCK {
					canMoveHorizontally = false
					break
				}
			}
		}
	}
	if !canMoveHorizontally {
		return false
	}

	minX := 0
	maxX := 0
	if dir == 1 {
		minX = len(chamber.rocks[0]) - 1
		maxX = -1
	} else {
		minX = 1
		maxX = len(chamber.rocks[0])
	}

	for y := chamber.fallingPieceLowestPoint; y < chamber.fallingPieceLowestPoint+shapeHeights[chamber.fallingPieceSeq]; y++ {
		for x := minX; x != maxX; x -= dir {
			if chamber.rocks[y][x] == CONTENTS_FALLING_ROCK {
				chamber.rocks[y][x] = 0
				chamber.rocks[y][x+dir] = CONTENTS_FALLING_ROCK
			}
		}
	}

	return true
}

func settleRock() {
	solidRockRow := -1
	for y := chamber.fallingPieceLowestPoint; y < chamber.fallingPieceLowestPoint+shapeHeights[chamber.fallingPieceSeq]; y++ {
		rowIsSolidRock := true
		for x := 0; x < len(chamber.rocks[y]); x++ {
			if chamber.rocks[y][x] == CONTENTS_FALLING_ROCK {
				chamber.rocks[y][x] = CONTENTS_SETTLED_ROCK
			} else if chamber.rocks[y][x] != CONTENTS_SETTLED_ROCK {
				rowIsSolidRock = false
			}
		}
		if rowIsSolidRock && y > solidRockRow {
			solidRockRow = y
		}
	}

	pieceHighest := chamber.fallingPieceLowestPoint + shapeHeights[chamber.fallingPieceSeq]
	if pieceHighest > chamber.highestSettledPoint {
		chamber.highestSettledPoint = pieceHighest
	}

	chamber.fallingPieceLowestPoint = -1
	chamber.fallingPieceSeq = -1

	if solidRockRow != -1 {
		trimChamber(solidRockRow)
	}
}

func trimChamber(solidRockRow int) {
	// fmt.Printf("Trimming chamber from row %d\n", solidRockRow)

	trashedRocks := chamber.rocks[:solidRockRow]
	savedRocks := chamber.rocks[solidRockRow:]
	chamber.rocks = append(savedRocks, trashedRocks...)

	for i := LONGEST_PERFORATED_STRETCH - 1; i >= solidRockRow; i-- {
		for j := 0; j <= CHAMBER_WIDTH-1; j++ {
			chamber.rocks[i][j] = 0
		}
	}

	chamber.heightBelowFLoor += solidRockRow
	chamber.highestSettledPoint -= solidRockRow
}

func applyGravity() bool {
	if chamber.fallingPieceSeq == -1 {
		return false
	}

	canFallDown := true
	for y := chamber.fallingPieceLowestPoint; y < chamber.fallingPieceLowestPoint+shapeHeights[chamber.fallingPieceSeq]; y++ {

		for x, contents := range chamber.rocks[y] {
			if contents == CONTENTS_FALLING_ROCK {
				if y == 0 || chamber.rocks[y-1][x] == CONTENTS_SETTLED_ROCK {
					canFallDown = false
					break
				}
			}
		}
	}

	if !canFallDown {
		settleRock()
		return false
	}

	for y := chamber.fallingPieceLowestPoint; y < chamber.fallingPieceLowestPoint+shapeHeights[chamber.fallingPieceSeq]; y++ {
		for x, contents := range chamber.rocks[y] {
			if contents == CONTENTS_FALLING_ROCK {
				chamber.rocks[y][x] = 0
				chamber.rocks[y-1][x] = CONTENTS_FALLING_ROCK
			}
		}
	}

	chamber.fallingPieceLowestPoint--
	return true
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
			contents := c.rocks[y][x]
			if contents == CONTENTS_SETTLED_ROCK {
				print("#")
			} else if contents == CONTENTS_FALLING_ROCK {
				print("@")
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

func doSimulation(totalRockCount int) {
	reportingCount := 5
	if totalRockCount > 10000 {
		reportingCount = 100000
	}
	reportingInterval := totalRockCount / reportingCount
	startTime := time.Now()

	makeChamber()

	rocksRemaining := totalRockCount
	shapeSeq := -1
	jetSeq := -1
	keyRockNumber := -1

	for rocksRemaining > 0 {
		if rocksRemaining%reportingInterval == 0 {
			elapsedSeconds := time.Since(startTime).Seconds()
			rocksCompleted := totalRockCount - rocksRemaining
			if rocksCompleted > 0 {
				remainingSeconds := (elapsedSeconds / float64(rocksCompleted)) * float64(rocksRemaining)

				fmt.Printf("%d rocks remaining, %d seconds remaining\n", rocksRemaining, int(remainingSeconds))
			}
		}

		rocksRemaining--
		shapeSeq++
		if shapeSeq >= len(shapes) {
			shapeSeq = 0
		}

		placeShape(shapeSeq)

		rockNumber := totalRockCount - rocksRemaining + 1
		// fmt.Printf("A new rock begins falling: #%d\n", rockNumber)
		// chamber.print()

		if rockNumber == keyRockNumber+1 {
			fmt.Printf("Lowest falling point: %d\n", chamber.fallingPieceLowestPoint)
			fmt.Printf("Highest settled point: %d\n", chamber.highestSettledPoint)
			panic("stop")
		}
		for {
			jetSeq++
			if jetSeq >= len(jetPattern) {
				jetSeq = 0
			}

			moved := applyJet(jetSeq)
			dir := "left"
			if jetPattern[jetSeq] == 1 {
				dir = "right"
			}

			if rockNumber == keyRockNumber {
				fmt.Printf("Jet of gas pushes rock %s", dir)
				if moved {
					println(":")
				} else {
					println(", but nothing happens:")
				}
				chamber.print()

			}

			movedDown := applyGravity()
			if rockNumber == keyRockNumber {
				print("Rock falls 1 unit")
				if movedDown {
					println(":")
				} else {
					println(", causing it to come to rest:")
				}
				chamber.print()
			}

			if !movedDown {
				break
			}
		}
	}

	println(chamber.highestSettledPoint + chamber.heightBelowFLoor)
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
