package main

import (
	"bufio"
	"os"
	"strings"
)

func indexOfUniqueStretch(message string, unique_len int) int {
	for cursor := unique_len - 1; cursor < len(message); cursor++ {
		chars := map[rune]bool{}
		for i := cursor - unique_len + 1; i <= cursor; i++ {
			chars[rune(message[i])] = true
		}

		if len(chars) == unique_len {
			return cursor + 1
		}
	}

	return -1
}

func part1(file *os.File) {
	sc := bufio.NewScanner(file)
	sc.Scan()

	println(indexOfUniqueStretch(sc.Text(), 4))
}

func part2(file *os.File) {
	sc := bufio.NewScanner(file)
	sc.Scan()

	println(indexOfUniqueStretch(sc.Text(), 14))
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
