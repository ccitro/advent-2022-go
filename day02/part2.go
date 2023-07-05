package main

import (
	"bufio"
	"os"
)

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

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

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
