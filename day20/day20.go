package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	data int
	seq  int
	next *Node
	prev *Node
}

type CyclicDoubleLinkedList struct {
	head     *Node
	capacity int
}

func (l *CyclicDoubleLinkedList) append(data int) {
	newNode := &Node{data, l.capacity, nil, nil}
	l.capacity++
	if l.head == nil {
		newNode.prev = newNode
		newNode.next = newNode
		l.head = newNode
		return
	}

	current := l.head
	for current.next != l.head {
		current = current.next
	}

	newNode.prev = current
	newNode.next = l.head
	current.next = newNode
	l.head.prev = newNode
}

func (l *CyclicDoubleLinkedList) print() {
	current := l.head
	for current.next != l.head {
		fmt.Printf("%d (seq %d), ", current.data, current.seq)
		current = current.next
	}
	fmt.Printf("%d (seq %d)\n", current.data, current.seq)

	fmt.Println()
}

func (l *CyclicDoubleLinkedList) findSeq(seq int) *Node {
	if seq < 0 || seq >= l.capacity {
		panic("Invalid sequence number")
	}

	current := l.head
	for current.seq != seq {
		current = current.next
	}
	return current
}

func (l *CyclicDoubleLinkedList) shiftNodeLeft(node *Node) {
	oneBefore := node.prev
	twoBefore := oneBefore.prev
	oneAfter := node.next

	twoBefore.next = node
	node.prev = twoBefore
	node.next = oneBefore
	oneBefore.prev = node
	oneBefore.next = oneAfter
	oneAfter.prev = oneBefore
}

func (l *CyclicDoubleLinkedList) shiftNodeRight(node *Node) {
	oneBefore := node.prev
	oneAfter := node.next
	twoAfter := oneAfter.next

	oneBefore.next = oneAfter
	oneAfter.prev = oneBefore
	oneAfter.next = node
	node.prev = oneAfter
	node.next = twoAfter
	twoAfter.prev = node
}

func (l *CyclicDoubleLinkedList) printAnswer() {
	l.print()

	current := l.head
	for {
		if current.data == 0 {
			break
		}
		current = current.next
	}

	for i := 0; i < 1000; i++ {
		current = current.next
	}
	plus1kval := current.data

	for i := 0; i < 1000; i++ {
		current = current.next
	}
	plus2kval := current.data

	for i := 0; i < 1000; i++ {
		current = current.next
	}
	plus3kval := current.data

	sum := plus1kval + plus2kval + plus3kval

	fmt.Printf("plus1kval: %d, plus2kval: %d, plus3kval: %d, sum: %d\n", plus1kval, plus2kval, plus3kval, sum)
}

func (l *CyclicDoubleLinkedList) mix() {
	for i := 0; i < l.capacity; i++ {
		node := l.findSeq(i)
		shiftAmount := node.data

		shiftAmount %= (l.capacity - 1)

		// these could be made more efficient by making a shiftNode method with a shiftAmount parameter
		// we're still fast enough with the current minimalist approach
		for shiftAmount > 0 {
			l.shiftNodeRight(node)
			shiftAmount--
		}
		for shiftAmount < 0 {
			l.shiftNodeLeft(node)
			shiftAmount++
		}
	}
}

var puzzleFile CyclicDoubleLinkedList

func loadPuzzle(file *os.File) {
	scanner := bufio.NewScanner(file)
	puzzleFile = CyclicDoubleLinkedList{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		intVal, _ := strconv.Atoi(line)
		puzzleFile.append(intVal)
	}
}

func part1() {
	println("Initial arrangement:")
	puzzleFile.print()
	puzzleFile.mix()
	puzzleFile.printAnswer()
}

func part2() {
	current := puzzleFile.head
	for {
		current.data *= 811589153
		current = current.next
		if current == puzzleFile.head {
			break
		}
	}

	println("Initial arrangement:")
	puzzleFile.print()
	for i := 0; i < 10; i++ {
		puzzleFile.mix()
	}
	puzzleFile.printAnswer()
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
