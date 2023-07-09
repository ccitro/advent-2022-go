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
	inspectCount    int
	items           []int
	operationChar   rune
	operationScalar int
	divisorTest     int
	successTarget   int
	failureTarget   int
}

func readMonkeys(file *os.File) []*monkey {
	scanner := bufio.NewScanner(file)
	monkeys := make([]*monkey, 0)
	seq := -1
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "Monkey") {
			continue
		}

		seq++
		m := monkey{id: seq, inspectCount: 0}
		monkeys = append(monkeys, &m)

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

func runRound(m *monkey, monkeys []*monkey, worryDivisor int) {
	fmt.Printf("Monkey %d:\n", m.id)

	for len(m.items) > 0 {
		v := m.items[0]
		m.items = m.items[1:]
		fmt.Printf("  Monkey inspects an item with worry level of %d\n", v)
		m.inspectCount++

		operationScalar := m.operationScalar
		if operationScalar == -1 {
			operationScalar = v
		}

		if m.operationChar == '+' {
			v += operationScalar
			fmt.Printf("    Worry level increases by %d to %d\n", operationScalar, v)
		} else {
			v *= operationScalar
			fmt.Printf("    Worry level is multiplied by %d to %d\n", operationScalar, v)
		}

		// sorry, I really should have put a mode in for this, but I'm lazy
		if worryDivisor == 3 {
			v /= worryDivisor // part 1
		} else {
			v %= worryDivisor // part 2
		}
		fmt.Printf("    Monkey gets bored with item. Worry level is divided by %d to %d\n", worryDivisor, v)

		isDivisible := v%m.divisorTest == 0
		target := -1
		if isDivisible {
			fmt.Printf("    Current worry level is divisible by %d\n", m.divisorTest)
			target = m.successTarget
		} else {
			fmt.Printf("    Current worry level is not divisible by %d\n", m.divisorTest)
			target = m.failureTarget
		}

		fmt.Printf("    Item with worry level %d is thrown to monkey %d\n", v, target)
		monkeys[target].items = append(monkeys[target].items, v)
	}
	println("")

}

func part1(file *os.File) {
	monkeys := readMonkeys(file)
	roundsRemaining := 20
	worryDivisor := 3

	for roundsRemaining > 0 {
		for _, m := range monkeys {
			runRound(m, monkeys, worryDivisor)
		}
		roundsRemaining--
	}
	for _, m := range monkeys {
		m.print()
	}

	inspectPlace1 := 0
	inspectPlace2 := 0
	for _, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times\n", m.id, m.inspectCount)
		if m.inspectCount > inspectPlace1 {
			inspectPlace2 = inspectPlace1
			inspectPlace1 = m.inspectCount
		} else if m.inspectCount > inspectPlace2 {
			inspectPlace2 = m.inspectCount
		}
	}

	fmt.Printf("The two monkeys who inspected the most items are %d and %d\n", inspectPlace1, inspectPlace2)
	monkeyBusiness := inspectPlace1 * inspectPlace2
	println(monkeyBusiness)
}

func part2(file *os.File) {
	monkeys := readMonkeys(file)
	roundsRemaining := 10000

	// uses the chinese remainder theorem to say that modular division
	// by the shared GCD of the divisors will have a unique remainder that
	// doesn't change the solution for each individual monkey
	worryDivisor := 1
	for _, m := range monkeys {
		worryDivisor *= m.divisorTest
	}

	for roundsRemaining > 0 {
		for _, m := range monkeys {
			runRound(m, monkeys, worryDivisor)
		}
		roundsRemaining--
	}
	for _, m := range monkeys {
		m.print()
	}

	inspectPlace1 := 0
	inspectPlace2 := 0
	for _, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times\n", m.id, m.inspectCount)
		if m.inspectCount > inspectPlace1 {
			inspectPlace2 = inspectPlace1
			inspectPlace1 = m.inspectCount
		} else if m.inspectCount > inspectPlace2 {
			inspectPlace2 = m.inspectCount
		}
	}

	fmt.Printf("The two monkeys who inspected the most items are %d and %d\n", inspectPlace1, inspectPlace2)
	monkeyBusiness := inspectPlace1 * inspectPlace2
	println(monkeyBusiness)

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
