package main

import (
	"bufio"
	"os"
	"strings"
)

func part1(file *os.File) {
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

// start part 2

type MatchScore int

const (
	WIN  MatchScore = 6
	DRAW            = 3
	LOSS            = 0
)

// opponent moves will come in as A, B, or C; remap to R, P, or S
var opponent_play_to_shape = map[string]string{"A": "R", "B": "P", "C": "S"}
var shape_to_score = map[string]int{"R": 1, "P": 2, "S": 3}

const OUTCOME_LOSE = "X"
const OUTCOME_DRAW = "Y"
const OUTCOME_WIN = "Z"

func determine_required_shape(opponent_shape string, desired_outcome string) string {
	if desired_outcome == OUTCOME_DRAW {
		return opponent_shape
	}

	if (desired_outcome == OUTCOME_WIN && opponent_shape == "R") ||
		(desired_outcome == OUTCOME_LOSE && opponent_shape == "S") {
		return "P"
	}

	if (desired_outcome == OUTCOME_WIN && opponent_shape == "P") ||
		(desired_outcome == OUTCOME_LOSE && opponent_shape == "R") {
		return "S"
	}

	return "R"
}

func score_match(my_shape string, op_shape string) MatchScore {
	if my_shape == op_shape {
		return DRAW
	}

	if my_shape == "R" && op_shape == "S" {
		return WIN
	}

	if my_shape == "P" && op_shape == "R" {
		return WIN
	}

	if my_shape == "S" && op_shape == "P" {
		return WIN
	}

	return LOSS
}

func part2(file *os.File) {
	score := 0

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		op_shape := opponent_play_to_shape[line[0:1]]
		desired_outcome := line[2:3]
		my_shape := determine_required_shape(op_shape, desired_outcome)
		match_score := score_match(my_shape, op_shape)
		shape_score := shape_to_score[my_shape]
		score += int(match_score) + shape_score
	}

	println(score)
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
