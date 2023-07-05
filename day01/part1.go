package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

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
