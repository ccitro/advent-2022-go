package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MonkeyPlan interface{}

type MonkeyPlanNumber int
type MonkeyPlanMath struct {
	left  string
	op    string
	right string
}

var monkeyPlans map[string]MonkeyPlan

func loadPuzzle(file *os.File) {
	monkeyPlans = make(map[string]MonkeyPlan)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			num, _ := strconv.Atoi(parts[1])
			id := parts[0][:len(parts[0])-1]
			monkeyPlans[id] = MonkeyPlanNumber(num)
			continue
		}

		if len(parts) == 4 {
			id := parts[0][:len(parts[0])-1]
			monkeyPlans[id] = MonkeyPlanMath{parts[1], parts[2], parts[3]}
			continue
		}

		panic("Unknown line: " + line)
	}
}

func evaluateFrom(id string) int {
	plan, ok := monkeyPlans[id]
	if !ok {
		panic("Unknown monkey: " + id)
	}

	switch plan := plan.(type) {
	case MonkeyPlanNumber:
		return int(plan)
	case MonkeyPlanMath:
		left := evaluateFrom(plan.left)
		right := evaluateFrom(plan.right)
		switch plan.op {
		case "+":
			return left + right
		case "*":
			return left * right
		case "-":
			return left - right
		case "/":
			return left / right

		default:
			panic("Unknown operator: " + plan.op)
		}

	default:
		panic("Unknown plan type")
	}
}

func part1() {
	val := evaluateFrom("root")
	println(val)
}

func part2() {
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
	loadPuzzle(file)
	method()
}
