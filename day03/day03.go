package main

import (
	"bufio"
	"os"
	"strings"
)

func getPriority(item_type string) int {
	ascii := int(item_type[0])
	if ascii > 96 {
		return ascii - 96
	}

	return ascii - 64 + 26
}

func part1(file *os.File) {
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

func part2(file *os.File) {
	// a slot in the array for each possible priority (a-z, A-Z)
	// the array is one larger than necessary so that the index matches the priority
	var rucksack [53]bool
	var shared_rucksack [53]int
	priority_sum := 0
	elf_seq := -1

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		elf_seq += 1

		for i := 0; i < len(line); i++ {
			priority := getPriority(line[i : i+1])
			// check if this priority has been seen for this elf
			if rucksack[priority] == false {
				// if it hasn't been seen, record it for this elf, and also increment the shared rucksack
				rucksack[priority] = true
				shared_rucksack[priority] += 1

				// if the shared rucksack has been incremented 3 times, then all
				// 3 elves have this priority, so add it to the sum
				if shared_rucksack[priority] == 3 {
					priority_sum += priority
					break
				}
			}
		}

		// reset rucksack for next line
		for i := range rucksack {
			rucksack[i] = false
		}

		// after the third elf in the group, reset shared rucksack
		if elf_seq == 2 {
			for i := range shared_rucksack {
				shared_rucksack[i] = 0
			}

			elf_seq = -1
		}
	}

	println(priority_sum)
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
