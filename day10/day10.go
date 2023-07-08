package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cpu struct {
	xreg int

	// value during nth cycle
	xreg_history []int
}

func buildCPU() cpu {
	return cpu{
		xreg:         1,
		xreg_history: []int{1},
	}
}

func (c *cpu) execute(instruction string) {
	if instruction == "" {
		return
	}

	parts := strings.Split(instruction, " ")
	opcode := parts[0]

	// noop - 1 cycle, no change
	if opcode == "noop" {
		c.xreg_history = append(c.xreg_history, c.xreg)
		return
	}

	// addx - add arg to xreg after two cycles
	amount, _ := strconv.Atoi(parts[1])
	c.xreg_history = append(c.xreg_history, c.xreg)
	c.xreg_history = append(c.xreg_history, c.xreg)
	c.xreg += amount
}

func part1(file *os.File) {
	cpu := buildCPU()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cpu.execute(scanner.Text())
	}

	// 20th cycle and every 40 cycles after that, up to 220
	key_cycles := []int{20, 60, 100, 140, 180, 220}
	total_strength := 0
	for _, v := range key_cycles {
		fmt.Printf("Cycle %d: %d\n", v, cpu.xreg_history[v])
		total_strength += cpu.xreg_history[v] * v
	}
	println(total_strength)
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
