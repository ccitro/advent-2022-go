package main

import (
	"bufio"
	"os"
)

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		println(line)
	}

}
