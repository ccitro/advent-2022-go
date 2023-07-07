package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type board struct {
	visited map[point]bool
	tail    point
	head    point
}

func (b *board) moveHeadOne(direction string) {
	switch direction {
	case "U":
		b.head.y++
	case "D":
		b.head.y--
	case "L":
		b.head.x--
	case "R":
		b.head.x++
	}

}

func (b *board) headAndTailAreTouching() bool {
	x_distance := b.head.x - b.tail.x
	if x_distance < 0 {
		x_distance = x_distance * -1
	}
	y_distance := b.head.y - b.tail.y
	if y_distance < 0 {
		y_distance = y_distance * -1
	}

	if x_distance == 0 && y_distance <= 1 {
		return true
	}

	if y_distance == 0 && x_distance <= 1 {
		return true
	}

	if x_distance == 1 && y_distance == 1 {
		return true
	}

	return false
}

func (b *board) dragTail() {
	if b.headAndTailAreTouching() {
		return
	}

	y_dir := 0
	if b.head.y > b.tail.y {
		y_dir = 1
	} else if b.head.y < b.tail.y {
		y_dir = -1
	}

	x_dir := 0
	if b.head.x > b.tail.x {
		x_dir = 1
	} else if b.head.x < b.tail.x {
		x_dir = -1
	}

	b.tail.x += x_dir
	b.tail.y += y_dir
}

func (b *board) move(m string) {
	direction := string(m[0])
	length, _ := strconv.Atoi(m[2:])

	for i := 0; i < length; i++ {
		b.moveHeadOne(direction)
		b.dragTail()
		b.print()
		b.visited[b.tail] = true
	}
}

func (b *board) print() {
	size := 6
	println("")
	for y := size; y >= 0; y-- {
		for x := 0; x <= size; x++ {
			if x == b.head.x && y == b.head.y {
				print("H")
			} else if x == b.tail.x && y == b.tail.y {
				print("T")
			} else if x == 0 && y == 0 {
				print("s")
			} else {
				print(".")
			}
		}
		println("")
	}
	println("")
}

func (b *board) printVisited() {
	size := 6
	println("")
	for y := size; y >= 0; y-- {
		for x := 0; x <= size; x++ {
			if x == 0 && y == 0 {
				print("s")
			} else if b.visited[point{x, y}] {
				print("X")
			} else {
				print(".")
			}
		}
		println("")
	}
	println("")

}

func part1(file *os.File) {
	board := board{visited: make(map[point]bool), tail: point{0, 0}, head: point{0, 0}}

	println("== Initial State ==")
	board.print()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		fmt.Printf("== %s ==\n", line)
		board.move(line)
	}

	board.printVisited()

	visit_count := 0
	for _, v := range board.visited {
		if v {
			visit_count++
		}
	}
	println(visit_count)
}

// start part2
// part 1 could be generalized to use part2 with a knot_count of 2,
// but I'm leaving the original solution in place for posterity

type part2board struct {
	visited map[point]bool
	knots   []point
}

func (b *part2board) moveHeadOne(direction string) {
	switch direction {
	case "U":
		b.knots[0].y++
	case "D":
		b.knots[0].y--
	case "L":
		b.knots[0].x--
	case "R":
		b.knots[0].x++
	}
}

func knotsAreTouching(a point, b point) bool {
	x_distance := a.x - b.x
	if x_distance < 0 {
		x_distance = x_distance * -1
	}
	y_distance := a.y - b.y
	if y_distance < 0 {
		y_distance = y_distance * -1
	}

	if x_distance == 0 && y_distance <= 1 {
		return true
	}

	if y_distance == 0 && x_distance <= 1 {
		return true
	}

	if x_distance == 1 && y_distance == 1 {
		return true
	}

	return false
}

func (b *part2board) dragTails() {
	for head_seq := 0; head_seq < len(b.knots)-1; head_seq++ {
		head := &b.knots[head_seq]
		tail := &b.knots[head_seq+1]

		if knotsAreTouching(*head, *tail) {
			continue
		}

		y_dir := 0
		if head.y > tail.y {
			y_dir = 1
		}
		if head.y < tail.y {
			y_dir = -1
		}

		x_dir := 0
		if head.x > tail.x {
			x_dir = 1
		}
		if head.x < tail.x {
			x_dir = -1
		}

		tail.x += x_dir
		tail.y += y_dir
	}
}

func (b *part2board) move(m string) {
	direction := string(m[0])
	length, _ := strconv.Atoi(m[2:])

	for i := 0; i < length; i++ {
		b.moveHeadOne(direction)
		b.dragTails()
		b.print()
		b.visited[b.knots[len(b.knots)-1]] = true
	}
}

func (b *part2board) getKnotNumberAtPoint(p point) int {
	for i, knot := range b.knots {
		if knot.x == p.x && knot.y == p.y {
			return i
		}
	}
	return -1
}

func (b *part2board) print() {
	size := 6
	println("")
	for y := size; y >= 0; y-- {
		for x := 0; x <= size; x++ {
			knot_number := b.getKnotNumberAtPoint(point{x, y})
			if knot_number == 0 {
				print("H")
			} else if knot_number > 0 {
				print(knot_number)
			} else if x == 0 && y == 0 {
				print("s")
			} else {
				print(".")
			}
		}
		println("")
	}
	println("")
}

func (b *part2board) printVisited() {
	size := 6
	println("")
	for y := size; y >= 0; y-- {
		for x := 0; x <= size; x++ {
			if x == 0 && y == 0 {
				print("s")
			} else if b.visited[point{x, y}] {
				print("X")
			} else {
				print(".")
			}
		}
		println("")
	}
	println("")
}

func createPart2Board(knot_count int) *part2board {
	visited := make(map[point]bool)
	knots := make([]point, knot_count)
	for i := 0; i < knot_count; i++ {
		knots[i] = point{0, 0}
	}
	board := part2board{visited: visited, knots: knots}
	return &board
}

func part2(file *os.File) {
	board := createPart2Board(10)
	println("== Initial State ==")
	board.print()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		fmt.Printf("== %s ==\n", line)
		board.move(line)
	}

	board.printVisited()

	visit_count := 0
	for _, v := range board.visited {
		if v {
			visit_count++
		}
	}
	println(visit_count)

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
