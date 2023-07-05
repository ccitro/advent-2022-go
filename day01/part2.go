package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

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
