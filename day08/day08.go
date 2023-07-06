package main

import (
	"bufio"
	"os"
	"strings"
)

type forest [][]int

var directions = [][]int{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

func readForest(file *os.File) *forest {
	sc := bufio.NewScanner(file)
	f := forest{}
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		r := make([]int, len(line))
		for i, c := range line {
			r[i] = int(c) - '0'
		}
		f = append(f, r)
	}

	return &f
}

func isVisible(f *forest, r int, c int) bool {
	h := len(*f)
	w := len((*f)[0])
	if r == 0 || c == 0 || r == h-1 || c == w-1 {
		return true
	}

	tree_height := (*f)[r][c]

	for _, d := range directions {
		x := r
		y := c

		all_shorter := true
		for {
			x += d[0]
			y += d[1]
			if x < 0 || y < 0 || x >= h || y >= w {
				break
			}

			if (*f)[x][y] >= tree_height {
				all_shorter = false
				break
			}
		}

		if all_shorter {
			return true
		}
	}

	return false
}

func calculateViewingDistance(f *forest, r int, c int, d []int) int {
	view_distance := 0
	source_tree_height := (*f)[r][c]

	for {
		r += d[0]
		c += d[1]
		if r < 0 || c < 0 || r >= len(*f) || c >= len((*f)[r]) {
			break
		}

		view_distance++
		loop_tree_height := (*f)[r][c]
		if loop_tree_height >= source_tree_height {
			break
		}
	}

	return view_distance
}

func calculateScenicScore(f *forest, r int, c int) int {
	score := 1

	for _, d := range directions {
		view_distance := calculateViewingDistance(f, r, c, d)
		score *= view_distance
	}

	return score
}

func part1(file *os.File) {
	forest := readForest(file)
	visible_count := 0
	for r := 0; r < len(*forest); r++ {
		for c := 0; c < len((*forest)[r]); c++ {
			if isVisible(forest, r, c) {
				visible_count++
			}
		}
	}

	println(visible_count)
}

func part2(file *os.File) {
	forest := readForest(file)
	highest_score := 0
	for r := 0; r < len(*forest); r++ {
		for c := 0; c < len((*forest)[r]); c++ {
			score := calculateScenicScore(forest, r, c)
			if score > highest_score {
				highest_score = score
			}
		}
	}

	println(highest_score)
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
