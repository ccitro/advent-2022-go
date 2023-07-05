package main

import (
	"bufio"
	"os"
)

func getPriority(item_type string) int {
	ascii := int(item_type[0])
	if ascii > 96 {
		return ascii - 96
	}

	return ascii - 64 + 26
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	var rucksack [53]int
	priority_sum := 0

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		rucksack_size := len(line)
		for i := 0; i < rucksack_size/2; i++ {
			priority := getPriority(line[i : i+1])
			rucksack[priority] = 1
		}

		for i := rucksack_size / 2; i < rucksack_size; i++ {
			priority := getPriority(line[i : i+1])
			if rucksack[priority] == 1 {
				priority_sum += priority
				break
			}
		}

		// reset rucksack for next line
		for i := 0; i < len(rucksack); i++ {
			rucksack[i] = 0
		}
	}

	println(priority_sum)
}
