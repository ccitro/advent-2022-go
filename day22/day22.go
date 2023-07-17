package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const BLOCK_VOID = 0
const BLOCK_WALL = 1
const BLOCK_OPEN = 2

const DIR_RIGHT = 0
const DIR_DOWN = 1
const DIR_LEFT = 2
const DIR_UP = 3

type Board [][]int

type PathNode interface{}
type PathNodeRotate rune
type PathNodeMove int
type Path []PathNode

type PuzzleState struct {
	x, y, dir, step int
}

type Pos struct {
	x, y int
}

var movements = map[int]Pos{
	DIR_RIGHT: {1, 0},
	DIR_DOWN:  {0, 1},
	DIR_LEFT:  {-1, 0},
	DIR_UP:    {0, -1},
}

var startX = -1
var startY = 0
var startDir = DIR_RIGHT
var board Board
var path Path
var state PuzzleState
var lastFacing map[Pos]int

func loadPuzzle(file *os.File) {
	maxWidth := -1
	height := 0

	// scan the file once to get the dimensions
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		height++
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
		if startX == -1 {
			startX = strings.Index(line, ".")
		}
	}

	board = make(Board, height)
	for i := range board {
		board[i] = make([]int, maxWidth)
	}
	path = make(Path, 0)
	lastFacing = make(map[Pos]int)

	// now loop again, this time loading the board
	y := 0
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(path) == 0 {
				// if we hit an empty line, then the next line is the path.  read it in
				scanner.Scan()
				line = scanner.Text()

				// parse the path char by char.  keep track of any digits we see, and add them as a
				// move length when we hit a non-digit
				digits := ""
				for _, v := range line {
					if v == 'R' || v == 'L' {
						if digits != "" {
							l, _ := strconv.Atoi(digits)
							path = append(path, PathNodeMove(l))
							digits = ""
						}
						path = append(path, PathNodeRotate(v))
					} else {
						digits += string(v)
					}
				}

				// if we have any digits left over, add them as a move length
				if digits != "" {
					l, _ := strconv.Atoi(digits)
					path = append(path, PathNodeMove(l))
				}
			}
			break
		}

		for x, v := range line {
			switch v {
			case '#':
				board[y][x] = BLOCK_WALL
			case '.':
				board[y][x] = BLOCK_OPEN
			}
		}

		y++
	}
}

func printPuzzleState() {
	fmt.Printf("Located at %d,%d facing %d on step %d\n", state.x, state.y, state.dir, state.step)
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[0]); x++ {
			facing, ok := lastFacing[Pos{x, y}]
			if ok {
				switch facing {
				case DIR_RIGHT:
					fmt.Print(">")
				case DIR_DOWN:
					fmt.Print("v")
				case DIR_LEFT:
					fmt.Print("<")
				case DIR_UP:
					fmt.Print("^")
				}
			} else {
				switch board[y][x] {
				case BLOCK_VOID:
					fmt.Print(" ")
				case BLOCK_WALL:
					fmt.Print("#")
				case BLOCK_OPEN:
					fmt.Print(".")
				}
			}
		}
		println()
	}
	println()
	for _, v := range path {
		switch v := v.(type) {
		case PathNodeRotate:
			fmt.Printf("%c", v)
		case PathNodeMove:
			fmt.Printf("%d", v)
		}
	}

	println()
	println()
}

func applyRotate(node PathNodeRotate) {
	switch node {
	case 'R':
		state.dir = (state.dir + 1) % 4
	case 'L':
		state.dir = (state.dir + 3) % 4
	default:
		panic("Unknown rotation")
	}
}

func wrapAround(pos Pos, dir int) Pos {
	newPos := Pos{pos.x, pos.y}
	for {
		if newPos.x >= len(board[0]) {
			newPos.x = 0
		}
		if newPos.y >= len(board) {
			newPos.y = 0
		}
		if newPos.x < 0 {
			newPos.x = len(board[0]) - 1
		}
		if newPos.y < 0 {
			newPos.y = len(board) - 1
		}

		if board[newPos.y][newPos.x] != BLOCK_VOID {
			return newPos
		}

		newPos.x += movements[dir].x
		newPos.y += movements[dir].y
	}
}

func applyMove(node PathNodeMove) {
	movesRemaining := int(node)
	for movesRemaining > 0 {
		lastFacing[Pos{state.x, state.y}] = state.dir
		movement := movements[state.dir]
		newPos := Pos{state.x + movement.x, state.y + movement.y}
		if newPos.x >= len(board[0]) || newPos.y >= len(board) || newPos.x < 0 || newPos.y < 0 || board[newPos.y][newPos.x] == BLOCK_VOID {
			newPos = wrapAround(newPos, state.dir)
		}

		if board[newPos.y][newPos.x] == BLOCK_WALL {
			break
		}

		state.x = newPos.x
		state.y = newPos.y
		movesRemaining--
	}
}

func applyPathNode(node PathNode) {
	switch node := node.(type) {
	case PathNodeRotate:
		applyRotate(node)
	case PathNodeMove:
		applyMove(node)
	default:
		panic("Unknown path node type")
	}
	state.step++
}

func part1() {
	state = PuzzleState{startX, startY, startDir, 0}
	lastFacing[Pos{startX, startY}] = startDir
	// printPuzzleState()

	for _, v := range path {
		applyPathNode(v)
		// printPuzzleState()
	}

	row := state.y + 1
	col := state.x + 1
	facing := state.dir

	password := 1000*row + 4*col + facing
	fmt.Printf("Ended at row=%d, col=%d, facing=%d, password=%d\n", row, col, facing, password)
	println(password)
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
