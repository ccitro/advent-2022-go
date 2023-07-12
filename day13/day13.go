package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// a packet is a list of values, each of which is either a number or a packet
type packetValue = interface{}
type packets = []packetValue

type pair struct {
	left  packets
	right packets
}

type part1Puzzle struct {
	pairs []pair
}

func parsePacket(line string) packets {
	// json decode line
	var p packets
	json.Unmarshal([]byte(line), &p)
	return p
}

func parseFileToPart1Puzzle(file *os.File) part1Puzzle {
	var puzzle part1Puzzle

	scanner := bufio.NewScanner(file)
	for {
		if !scanner.Scan() {
			break
		}
		line1 := scanner.Text()
		if len(line1) == 0 {
			continue
		}

		if !scanner.Scan() {
			break
		}
		line2 := scanner.Text()

		line1packet := parsePacket(line1)
		line2packet := parsePacket(line2)
		puzzle.pairs = append(puzzle.pairs, pair{line1packet, line2packet})
	}

	return puzzle
}

func compareValues(left packetValue, right packetValue) int {
	// If the lists are the same length and no comparison makes a decision about the order, continue checking the next part of the input.
	if left == nil && right == nil {
		return 0
	}

	// If the left list runs out of items first, the inputs are in the right order
	if left == nil {
		return -1
	}

	// If the right list runs out of items first, the inputs are not in the right order
	if right == nil {
		return 1
	}

	// If both values are numbers, the lower number should come first
	leftNumber, leftIsNumber := left.(float64)
	rightNumber, rightIsNumber := right.(float64)

	if leftIsNumber && rightIsNumber {
		if leftNumber < rightNumber {
			return -1
		}
		if leftNumber > rightNumber {
			return 1
		}
		return 0
	}

	// if one value is an number and the other is a list, promote the number to a list
	var leftList packets
	var rightList packets

	if leftIsNumber {
		leftList = packets{left}
		rightList = right.(packets)
	} else if rightIsNumber {
		rightList = packets{right}
		leftList = left.(packets)
	} else {
		leftList = left.(packets)
		rightList = right.(packets)
	}

	i := 0
	for {
		var listLeftValue packetValue
		var listRightValue packetValue
		if i < len(leftList) {
			listLeftValue = leftList[i]
		}
		if i < len(rightList) {
			listRightValue = rightList[i]
		}

		if listLeftValue == nil && listRightValue == nil {
			return 0
		}

		comp := compareValues(listLeftValue, listRightValue)
		if comp != 0 {
			return comp
		}
		i++
	}
}

func isPairInOrder(p pair) bool {
	comp := compareValues(p.left, p.right)
	return comp < 0
}

func (p *part1Puzzle) print() {
	for _, v := range p.pairs {
		fmt.Printf("%v\n", v.left)
		fmt.Printf("%v\n", v.right)
		println()
	}
}

func part1(file *os.File) {
	puzzle := parseFileToPart1Puzzle(file)

	sum := 0
	for i, v := range puzzle.pairs {
		seq := i + 1
		inRightOrder := isPairInOrder(v)
		fmt.Printf("Pair %d: %v\n", seq, inRightOrder)
		if inRightOrder {
			sum += seq
		}
	}
	println(sum)
}

type part2Puzzle struct {
	packets packets
}

func parseFileToPart2Puzzle(file *os.File) part2Puzzle {
	var puzzle part2Puzzle

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		packet := parsePacket(line)
		puzzle.packets = append(puzzle.packets, packet)
	}

	return puzzle
}

// func (p packets) Len() int {
// 	return len(p)
// }

// func (p packets) Less(i, j int) bool {
// 	return compareValues(p[i], p[j]) < 0
// }

// func (p packets) Swap(i, j int) {
// 	p[i], p[j] = p[j], p[i]
// }

func sortPackets(p packets) []packetValue {
	// terribly inefficient, but I didn't spend the time to figure out how to implement the sort interface for {}interface
	sorted := make(packets, len(p))
	copy(sorted, p)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if compareValues(sorted[i], sorted[j]) > 0 {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func part2(file *os.File) {
	puzzle := parseFileToPart2Puzzle(file)

	dividerPacketsJson := []string{
		"[[2]]",
		"[[6]]",
	}

	for _, v := range dividerPacketsJson {
		p := parsePacket(v)
		puzzle.packets = append(puzzle.packets, p)
	}

	// sorted := sortPackets(puzzle.packets)
	sorted := sortPackets(puzzle.packets)
	dividerPacket0Index := -1
	dividerPacket1Index := -1
	for i, v := range sorted {
		// I'm not familiar enough with go to know if there's a better way to compare two "packets", since
		// they're interface{} types.  converting it back to strings/json works well enough, even if its not ideal.
		vJson := fmt.Sprintf("%v", v)
		if vJson == dividerPacketsJson[0] {
			dividerPacket0Index = i + 1
		} else if vJson == dividerPacketsJson[1] {
			dividerPacket1Index = i + 1
		}
	}

	println(dividerPacket0Index * dividerPacket1Index)
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
