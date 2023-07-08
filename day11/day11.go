package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type monkey struct {
	id              int
	items           []int
	operationChar   rune
	operationScalar int
	divisorTest     int
	successTarget   int
	failureTarget   int
}

func readMonkeys(file *os.File) []monkey {
	scanner := bufio.NewScanner(file)
	monkeys := make([]monkey, 0)
	seq := -1
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "Monkey") {
			continue
		}

		seq++
		monkeys = append(monkeys, monkey{id: seq})

		scanner.Scan()
		line = scanner.Text() // Starting items
		itemsList := strings.Split(line[18:], ", ")
		monkeys[seq].items = make([]int, len(itemsList))
		for i, v := range itemsList {
			monkeys[seq].items[i], _ = strconv.Atoi(v)
		}

		scanner.Scan()
		line = scanner.Text() // Operation
		monkeys[seq].operationChar = rune(line[23])
		operationText := line[25:]
		if operationText == "old" {
			monkeys[seq].operationScalar = -1
		} else {
			monkeys[seq].operationScalar, _ = strconv.Atoi(operationText)
		}

		scanner.Scan()
		line = scanner.Text() // Test
		monkeys[seq].divisorTest, _ = strconv.Atoi(line[21:])

		scanner.Scan()
		line = scanner.Text() // If true
		monkeys[seq].successTarget, _ = strconv.Atoi(line[29:])

		scanner.Scan()
		line = scanner.Text() // If false
		monkeys[seq].failureTarget, _ = strconv.Atoi(line[30:])
	}

	return monkeys
}

func (m *monkey) print() {
	fmt.Printf("Monkey %d:\n", m.id)
	fmt.Printf("  Starting items: %v\n", m.items)

	operationScalarText := strconv.Itoa(m.operationScalar)
	if m.operationScalar == -1 {
		operationScalarText = "old"
	}
	fmt.Printf("  Operation: new = old %c %s\n", m.operationChar, operationScalarText)

	fmt.Printf("  Test: divisible by %d\n", m.divisorTest)
	fmt.Printf("    If true: throw to monkey %d\n", m.successTarget)
	fmt.Printf("    If false: throw to monkey %d\n", m.failureTarget)
	println("")
}

func part1(file *os.File) {
	monkeys := readMonkeys(file)
	for _, v := range monkeys {
		v.print()
	}
}

func part2(file *os.File) {
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
