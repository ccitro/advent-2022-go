package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MAX_LEN = 20

var lavaDroplet [MAX_LEN][MAX_LEN][MAX_LEN]bool

func loadPuzzle(file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		x := 0
		y := 0
		z := 0

		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		lavaDroplet[x][y][z] = true
	}
}

func part1() {
	exposedSides := 0
	for x := 0; x < MAX_LEN; x++ {
		for y := 0; y < MAX_LEN; y++ {
			for z := 0; z < MAX_LEN; z++ {
				if !lavaDroplet[x][y][z] {
					continue
				}

				if x == 0 || !lavaDroplet[x-1][y][z] {
					exposedSides++
				}
				if x == MAX_LEN-1 || !lavaDroplet[x+1][y][z] {
					exposedSides++
				}
				if y == 0 || !lavaDroplet[x][y-1][z] {
					exposedSides++
				}
				if y == MAX_LEN-1 || !lavaDroplet[x][y+1][z] {
					exposedSides++
				}
				if z == 0 || !lavaDroplet[x][y][z-1] {
					exposedSides++
				}
				if z == MAX_LEN-1 || !lavaDroplet[x][y][z+1] {
					exposedSides++
				}
			}
		}
	}

	println(exposedSides)
}

func part2() {
	dirs := [6][3]int{
		{-1, 0, 0},
		{1, 0, 0},
		{0, -1, 0},
		{0, 1, 0},
		{0, 0, -1},
		{0, 0, 1},
	}

	// repeatedly loop over every square, and propogate the "exposedness" state of every air square, where
	// a square is exposed if it is on an edge, or if it is adjacent to an exposed air square
	// this is really bad for time complexity, but with n=20 it doesn't matter
	exposedAirSquares := [MAX_LEN][MAX_LEN][MAX_LEN]bool{}
	changeMade := true
	for changeMade {
		changeMade = false
		for x := 0; x < MAX_LEN; x++ {
			for y := 0; y < MAX_LEN; y++ {
				for z := 0; z < MAX_LEN; z++ {
					if exposedAirSquares[x][y][z] {
						// fmt.Printf("Square at %d,%d,%d is already exposed\n", x, y, z)
						continue
					}
					if lavaDroplet[x][y][z] {
						// fmt.Printf("Square at %d,%d,%d is lava\n", x, y, z)
						continue
					}

					if x == 0 || y == 0 || z == 0 || x == MAX_LEN-1 || y == MAX_LEN-1 || z == MAX_LEN-1 {
						// fmt.Printf("Square at %d,%d,%d is exposed due to edge\n", x, y, z)
						exposedAirSquares[x][y][z] = true
						changeMade = true
						continue
					}

					for _, dir := range dirs {
						xx := x + dir[0]
						yy := y + dir[1]
						zz := z + dir[2]
						if exposedAirSquares[xx][yy][zz] {
							// fmt.Printf("Square at %d,%d,%d is exposed due to adjacent air\n", x, y, z)
							exposedAirSquares[x][y][z] = true
							changeMade = true
							break
						}
					}
				}
			}
		}
	}

	// fmt.Printf("Exposed air check at 2,2,5: %v\n", exposedAirSquares[2][2][5])

	// now repeat the same type of loop as part1, but only count squares that are adjacent to exposed air
	exteriorSides := 0
	for x := 0; x < MAX_LEN; x++ {
		for y := 0; y < MAX_LEN; y++ {
			for z := 0; z < MAX_LEN; z++ {
				if !lavaDroplet[x][y][z] {
					continue
				}
				for _, dir := range dirs {
					xx := x + dir[0]
					yy := y + dir[1]
					zz := z + dir[2]
					if xx < 0 || yy < 0 || zz < 0 || xx >= MAX_LEN || yy >= MAX_LEN || zz >= MAX_LEN {
						exteriorSides++
						continue
					}

					if exposedAirSquares[xx][yy][zz] {
						exteriorSides++
						continue
					}
				}
			}
		}
	}

	println(exteriorSides)
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
