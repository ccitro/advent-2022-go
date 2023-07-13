package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

func part2(file *os.File) {
	sensorData := readSensorData(file)
	searchRange := 4000000
	// searchRange = 20
	var blockingSensorData *SensorData

	lastTs := time.Now()
	for x := 0; x < searchRange; x++ {
		for y := 0; y < searchRange; y++ {
			// fmt.Printf("Checking x=%d, y=%d\n", x, y)
			if (x%1000000) == 10 && y == 0 {
				n := time.Now()
				elapsed := n.Sub(lastTs).Seconds()
				percentComplete := float64(x) / float64(searchRange)
				estimatedTimeRemaining := elapsed / percentComplete
				fmt.Printf("x=%d, y=%d, percent complete: %f, estimated time remaining: %f\n", x, y, percentComplete, estimatedTimeRemaining)
			}

			blockingSensorData = nil

			for _, v := range *sensorData {
				if (abs(v.x-x) + abs(v.y-y)) <= v.sensorRange {
					blockingSensorData = &v
					break
				}
			}

			if blockingSensorData != nil {
				// we have just entered a diamond/triangle that is blocked by a sensor
				// figure out the furthest Y value of that shape along this X row

				xDistanceToSensor := abs(blockingSensorData.x - x)

				// manhattan distance to the sensor means that if we know how far the x distance is,
				// we can figure out the y reach of the sensor
				yReach := blockingSensorData.sensorRange - xDistanceToSensor

				// calculate the max y value of the region this sensor covers, using the reach we calculated
				sensorMaxY := blockingSensorData.y + yReach

				// fmt.Printf("I am at %d, %d and I am blocked by a sensor at %d, %d that has a reach of %d.  The x distance to the sensor is %d, the y reach is %d, and the sensor max y is %d\n", x, y, blockingSensorData.x, blockingSensorData.y, blockingSensorData.sensorRange, xDistanceToSensor, yReach, sensorMaxY)
				if y < sensorMaxY {
					y = sensorMaxY
				}

				continue
			}

			tuningFrequency := 4000000*x + y
			fmt.Printf("Found a spot at x=%d, y=%d, tuning frequency: %d\n", x, y, tuningFrequency)
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
