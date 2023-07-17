package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

var strValues = map[string]int{
	"2": 2,
	"1": 1,
	"0": 0,
	"-": -1,
	"=": -2,
}

var valueStrs = map[int]string{
	2:  "2",
	1:  "1",
	0:  "0",
	-1: "-",
	-2: "=",
}

var fuelRequirements []string

func loadPuzzle(file *os.File) {
	fuelRequirements = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fuelRequirements = append(fuelRequirements, line)
	}
}

func snafuToDecimal(s string) int {
	total := 0
	for i := len(s) - 1; i >= 0; i-- {
		place := len(s) - i - 1
		placeMulti := int(math.Pow(5, float64(place)))
		strValue := strValues[string(s[i])]
		placeValue := strValue * placeMulti
		total += placeValue
	}

	return total
}

func reverse(a []string) []string {
	b := make([]string, len(a))
	for i := len(a) - 1; i >= 0; i-- {
		b[len(a)-i-1] = a[i]
	}

	return b
}

func addSnafu(a, b string) string {
	digits1 := strings.Split(a, "")
	digits2 := strings.Split(b, "")

	shorter := reverse(digits1)
	longer := reverse(digits2)
	if len(digits1) > len(digits2) {
		shorter, longer = longer, shorter
	}

	ans := make([]int, len(longer)+1)
	for i := 0; i < len(longer); i++ {
		ans[i] += strValues[longer[i]]
		if i < len(shorter) {
			ans[i] += strValues[shorter[i]]
		}
		if ans[i] > 2 {
			ans[i] -= 5
			ans[i+1]++
		} else if ans[i] < -2 {
			ans[i] += 5
			ans[i+1]--
		}
	}

	snafu := ""
	for i := len(ans) - 1; i >= 0; i-- {
		snafu += valueStrs[ans[i]]
	}
	snafu = strings.TrimLeft(snafu, "0")
	if snafu == "" {
		snafu = "0"
	}

	return snafu
}

func part1() {
	// initially I added in decimal, but convering from dec to snafu is a pain
	total := 0
	for _, v := range fuelRequirements {
		total += snafuToDecimal(v)
	}
	fmt.Printf("Total req: %d\n", total)

	// so just do the addition in snafu
	snafuTotal := "0"
	for _, v := range fuelRequirements {
		snafuTotal = addSnafu(snafuTotal, v)
	}
	fmt.Printf("Snafu total: %s\n", snafuTotal)
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
