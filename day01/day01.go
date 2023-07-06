package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func part1(file *os.File) {
	max := 0
	current := 0

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			if current > max {
				max = current
			}
			current = 0
		} else {
			calories, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			current += calories
		}
	}

	println(max)
}

func part2(file *os.File) {
	maxes := []int{0, 0, 0}
	current := 0

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		calories, err := strconv.Atoi(line)
		if err == nil {
			current += calories
			continue
		}

		if current > maxes[0] {
			maxes[2] = maxes[1]
			maxes[1] = maxes[0]
			maxes[0] = current
		} else if current > maxes[1] {
			maxes[2] = maxes[1]
			maxes[1] = current
		} else if current > maxes[2] {
			maxes[2] = current
		}

		current = 0
	}

	println(maxes[0] + maxes[1] + maxes[2])
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
