package main

import (
	"bufio"
	"fmt"
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

func planInvolvesHuman(id string) bool {
	if id == "humn" {
		return true
	}

	plan, ok := monkeyPlans[id]
	if !ok {
		panic("Unknown monkey: " + id)
	}

	switch plan := plan.(type) {
	case MonkeyPlanNumber:
		return false
	case MonkeyPlanMath:
		return planInvolvesHuman(plan.left) || planInvolvesHuman(plan.right)
	default:
		panic("Unknown plan type")
	}
}

func getExpression(id string) string {
	if id == "humn" {
		return "x"
	}
	plan := monkeyPlans[id]
	switch plan := plan.(type) {
	case MonkeyPlanNumber:
		return fmt.Sprintf("%d", plan)
	case MonkeyPlanMath:
		left := ""
		right := ""
		if planInvolvesHuman(plan.left) {
			left = getExpression(plan.left)
			right = fmt.Sprintf("%d", evaluateFrom(plan.right))
		} else {
			left = fmt.Sprintf("%d", evaluateFrom(plan.left))
			right = getExpression(plan.right)
		}

		if plan.op == "*" || plan.op == "/" {
			return fmt.Sprintf("%s%s%s", left, plan.op, right)
		}
		return fmt.Sprintf("(%s%s%s)", left, plan.op, right)
	default:
		panic("Unknown plan type")
	}
}

func part2() {
	rootMonkeyPlan := monkeyPlans["root"].(MonkeyPlanMath)

	humanTree := ""
	monkeyTreeValue := -1

	if planInvolvesHuman(rootMonkeyPlan.left) {
		humanTree = rootMonkeyPlan.left
		monkeyTreeValue = evaluateFrom(rootMonkeyPlan.right)
	} else {
		humanTree = rootMonkeyPlan.right
		monkeyTreeValue = evaluateFrom(rootMonkeyPlan.left)
	}

	fmt.Printf("Need to make tree starting at %s equal %d\n", humanTree, monkeyTreeValue)

	textDescription := getExpression(humanTree)
	fmt.Printf("%d=%s\n", monkeyTreeValue, textDescription)

	// the next step is to solve the equation above for x
	// i used an external tool to do this
	// @todo consider implementing this in go, as an equation solver or some sort of tree rebalancer?
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
