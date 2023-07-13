package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SensorData struct {
	x           int
	y           int
	beaconX     int
	beaconY     int
	sensorRange int
}

func readSensorData(file *os.File) *[]SensorData {
	var sensorData []SensorData
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		x := 0
		y := 0
		beaconX := 0
		beaconY := 0
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &x, &y, &beaconX, &beaconY)
		xDistance := abs(x - beaconX)
		yDistance := abs(y - beaconY)
		sensorRange := xDistance + yDistance
		sensorData = append(sensorData, SensorData{x, y, beaconX, beaconY, sensorRange})
	}

	return &sensorData
}

func (s *SensorData) print() {
	fmt.Printf("Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d\n", s.x, s.y, s.beaconX, s.beaconY)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isBeaconAt(x, y int, sensorData *[]SensorData) bool {
	// this could be stored in a map for faster lookup, but the array isn't large so this works fine
	for _, v := range *sensorData {
		if v.beaconX == x && v.beaconY == y {
			return true
		}
	}

	return false
}

func coordsCanHoldBeacon(x, y int, sensorData *[]SensorData) bool {
	if isBeaconAt(x, y, sensorData) {
		return true
	}

	for _, v := range *sensorData {
		xDistanceToSensor := abs(v.x - x)
		yDistanceToSensor := abs(v.y - y)
		distanceToSensor := xDistanceToSensor + yDistanceToSensor
		if distanceToSensor <= v.sensorRange {
			return false
		}
	}

	return true
}

func part1(file *os.File) {
	sensorData := readSensorData(file)
	minX := 99999
	maxX := -99999

	for _, v := range *sensorData {
		sensorMinX := v.x - v.sensorRange
		sensorMaxX := v.x + v.sensorRange
		if sensorMinX < minX {
			minX = sensorMinX
		}
		if sensorMaxX > maxX {
			maxX = sensorMaxX
		}
	}
	fmt.Printf("minX: %d, maxX: %d\n", minX, maxX)

	row := 2000000
	blockedPosCount := 0
	for i := minX; i <= maxX; i++ {
		if !coordsCanHoldBeacon(i, row, sensorData) {
			blockedPosCount++
		}
	}

	println(blockedPosCount)
}

type pos struct {
	x int
	y int
}

func posCanHoldBeacon(p pos, sensorData *[]SensorData) bool {
	for _, v := range *sensorData {
		xDistanceToSensor := abs(v.x - p.x)
		yDistanceToSensor := abs(v.y - p.y)
		distanceToSensor := xDistanceToSensor + yDistanceToSensor
		if distanceToSensor <= v.sensorRange {
			return false
		}
	}

	return true
}

func part2(file *os.File) {
	sensorData := readSensorData(file)
	searchRange := 4000000

	beaconLocations := make(map[pos]bool)
	for _, v := range *sensorData {
		beaconLocations[pos{v.beaconX, v.beaconY}] = true
	}

	sensorLocations := make(map[pos]bool)
	for _, v := range *sensorData {
		sensorLocations[pos{v.x, v.y}] = true
	}

	for x := 0; x < searchRange; x++ {
		for y := 0; y < searchRange; y++ {
			if (x%10) == 0 && y == 0 {
				fmt.Printf("x=%d, y=%d\n", x, y)
			}
			p := pos{x, y}
			if (!beaconLocations[p]) && (!sensorLocations[p]) && posCanHoldBeacon(p, sensorData) {
				tuning := 4000000*x + y
				fmt.Printf("Beacon is possible at x=%d, y=%d, tuning=%d\n", x, y, tuning)
			}
		}
	}
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
