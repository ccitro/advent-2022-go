package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type assignment struct {
	start int
	end   int
}

func assignment_contains(a assignment, b assignment) bool {
	return a.start <= b.start && a.end >= b.end
}

func assignment_overlaps(a assignment, b assignment) bool {
	return (a.start >= b.start && a.start <= b.end) || (a.end >= b.start && a.end <= b.end)
}

func parse_assignment(s string) assignment {
	parts := strings.Split(s, "-")
	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])
	return assignment{start, end}
}

func part1(file *os.File) {
	overlapping := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		assignments := strings.Split(line, ",")
		assignment1 := parse_assignment(assignments[0])
		assignment2 := parse_assignment(assignments[1])
		if assignment_contains(assignment1, assignment2) || assignment_contains(assignment2, assignment1) {
			overlapping++
		}
	}

	println(overlapping)
}

func part2(file *os.File) {
	overlapping := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		assignments := strings.Split(line, ",")
		assignment1 := parse_assignment(assignments[0])
		assignment2 := parse_assignment(assignments[1])
		if assignment_overlaps(assignment1, assignment2) || assignment_overlaps(assignment2, assignment1) {
			overlapping++
		}
	}

	println(overlapping)
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
