package main

import (
	"bufio"
	"os"
)

func main() {
	shape_scores := map[string]int{"X": 1, "Y": 2, "Z": 3}
	const SCORE_WIN = 6
	const SCORE_DRAW = 3
	const SCORE_LOSS = 0

	outcomes := map[string]int{
		"A X": SCORE_DRAW,
		"A Y": SCORE_WIN,
		"A Z": SCORE_LOSS,
		"B X": SCORE_LOSS,
		"B Y": SCORE_DRAW,
		"B Z": SCORE_WIN,
		"C X": SCORE_WIN,
		"C Y": SCORE_LOSS,
		"C Z": SCORE_DRAW,
	}

	file, _ := os.Open("input.txt")
	defer file.Close()

	score := 0

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		my_shape := line[2:3]
		score += shape_scores[my_shape] + outcomes[line]
	}

	println(score)
}
