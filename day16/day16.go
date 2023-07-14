package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	name    string
	flow    int
	tunnels []string
}

type Graph = map[string]*Node

type Route struct {
	flow  int
	nodes []string
}

var graph Graph
var distances map[string]map[string]int
var usefulValves []string

func readPuzzleGraph(file *os.File) {
	graph = Graph{}
	usefulValves = []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		name := ""
		flow := 0
		fmt.Sscanf(line, "Valve %s has flow rate=%d; tunnels lead to valves", &name, &flow)
		suffix := strings.Split(line, "to valv")[1]
		suffix = strings.Replace(suffix, ",", "", -1)
		tunnels := strings.Split(suffix, " ")[1:]

		node := Node{name: name, flow: flow, tunnels: tunnels}
		if node.flow > 0 {
			usefulValves = append(usefulValves, name)
		}
		graph[name] = &node
	}
}

func floydWarshall() {
	distances = make(map[string]map[string]int)

	for src, srcNode := range graph {
		distances[src] = make(map[string]int)
		for dest := range graph {
			if src == dest {
				distances[src][dest] = 0
			} else {
				distances[src][dest] = 1e9
			}
		}
		for _, tunnel := range srcNode.tunnels {
			distances[src][tunnel] = 1
		}
	}

	for k := range graph {
		for i := range graph {
			for j := range graph {
				if distances[i][j] > distances[i][k]+distances[k][j] {
					distances[i][j] = distances[i][k] + distances[k][j]
				}
			}
		}
	}
}

func searchRoutes(start string, time int, route Route, visited map[string]bool) []Route {
	routes := []Route{route}

	for _, valve := range usefulValves {
		newTime := time - distances[start][valve] - 1
		if visited[valve] || newTime < 0 {
			continue
		}

		newVisited := make(map[string]bool)
		for k, v := range visited {
			newVisited[k] = v
		}
		newVisited[valve] = true

		newRoute := Route{}
		newRoute.flow = route.flow + graph[valve].flow*newTime
		newRoute.nodes = make([]string, len(route.nodes))
		copy(newRoute.nodes, route.nodes)
		newRoute.nodes = append(newRoute.nodes, valve)

		routes = append(routes, searchRoutes(valve, newTime, newRoute, newVisited)...)
	}

	return routes
}

func part1(file *os.File) {
	readPuzzleGraph(file)
	start := "AA"
	duration := 30

	initialRoute := Route{flow: 0, nodes: []string{start}}
	floydWarshall()
	visited := make(map[string]bool)

	routes := searchRoutes(start, duration, initialRoute, visited)
	bestRoute := routes[0]
	for _, route := range routes {
		if route.flow > bestRoute.flow {
			bestRoute = route
		}
	}

	fmt.Printf("Best route: %v\n", bestRoute)
}

func allDifferentNodes(nodes1 []string, nodes2 []string) bool {
	for _, node1 := range nodes1 {
		for _, node2 := range nodes2 {
			if node1 == node2 && node1 != "AA" {
				return false
			}
		}
	}
	return true
}

func part2(file *os.File) {
	readPuzzleGraph(file)
	start := "AA"
	duration := 26

	initialRoute := Route{flow: 0, nodes: []string{start}}
	floydWarshall()

	visited := make(map[string]bool)
	routes := searchRoutes(start, duration, initialRoute, visited)

	max := 0
	for _, myRoute := range routes {
		if len(myRoute.nodes) > 0 {
			for _, elephantRoute := range routes {
				totalFlow := myRoute.flow + elephantRoute.flow
				if totalFlow <= max {
					continue
				}

				if allDifferentNodes(myRoute.nodes, elephantRoute.nodes) {
					max = totalFlow
				}
			}
		}
	}

	fmt.Printf("Max flow: %d\n", max)
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
