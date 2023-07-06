package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type move struct {
	qty    int
	source int
	dest   int
}

type board [10][]byte

type game struct {
	board board
	moves []move
}

func getResult(b *board) string {
	result := ""
	for i := range b {
		t := b[i]
		if len(t) > 0 {
			result += string(t[0])
		}
	}
	return result
}

func readMove(line string) move {
	m := move{}
	parts := strings.Split(line, " ")
	m.qty, _ = strconv.Atoi(parts[1])
	m.source, _ = strconv.Atoi(parts[3])
	m.dest, _ = strconv.Atoi(parts[5])

	// 0-index for easier slice access
	m.source--
	m.dest--

	return m
}

func readGame(file *os.File) *game {
	sc := bufio.NewScanner(file)
	g := game{}
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if line[1] == '1' {
			continue
		}

		if strings.HasPrefix(line, "move") {
			g.moves = append(g.moves, readMove(line))
			continue
		}

		for i := 1; i < len(line); i += 4 {
			r := line[i]
			if r != ' ' {
				tower := (i - 1) / 4
				g.board[tower] = append(g.board[tower], r)
			}
		}
	}

	return &g
}

func executePart2Move(b *board, m move) {
	blocks_to_move := append([]byte{}, b[m.source][:m.qty]...)
	b[m.source] = b[m.source][m.qty:]
	b[m.dest] = append(blocks_to_move, b[m.dest]...)
}

func executePart1Move(b *board, m move) {
	// part 1 can only move blocks 1 at a time, which is equivalent to moving all the blocks at once in reverse order
	// pull the blocks to move out of the source by tower, one by one, in reverse order
	blocks_to_move := make([]byte, m.qty)
	for i := m.qty - 1; i >= 0; i-- {
		dest_i := m.qty - i - 1
		blocks_to_move[dest_i] = b[m.source][i]
	}

	b[m.source] = b[m.source][m.qty:]
	b[m.dest] = append(blocks_to_move, b[m.dest]...)
}

func part1(file *os.File) {
	game := readGame(file)
	for _, move := range game.moves {
		executePart1Move(&game.board, move)
	}
	println(getResult(&game.board))
}

func part2(file *os.File) {
	game := readGame(file)
	for _, move := range game.moves {
		executePart2Move(&game.board, move)
	}
	println(getResult(&game.board))
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
