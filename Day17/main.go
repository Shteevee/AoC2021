package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x int
	y int
}

type TargetBounds struct {
	yLower int
	yUpper int
	xLower int
	xUpper int
}

func calculateTriangle(x int) int {
	return (x * (x + 1)) / 2
}

func parseArea(scanner *bufio.Scanner) TargetBounds {
	scanner.Scan()
	coords := strings.Split(strings.TrimPrefix(scanner.Text(), "target area: "), ",")
	xRangeString := strings.Split(strings.TrimPrefix(coords[0], "x="), "..")
	yRangeString := strings.Split(strings.TrimPrefix(coords[1], " y="), "..")

	xLower, _ := strconv.Atoi(xRangeString[0])
	xUpper, _ := strconv.Atoi(xRangeString[1])
	yLower, _ := strconv.Atoi(yRangeString[0])
	yUpper, _ := strconv.Atoi(yRangeString[1])

	return TargetBounds{
		yLower: yLower,
		yUpper: yUpper,
		xLower: xLower,
		xUpper: xUpper,
	}
}

func findMinX(xLower int) int {
	for i := 0; i < xLower; i++ {
		maxDistance := calculateTriangle(i)
		if maxDistance >= xLower {
			return i
		}
	}

	return xLower
}

func trajectoryHitsTarget(x int, y int, targetBounds TargetBounds) bool {
	currentPoint := Point{x: 0, y: 0}
	for currentPoint.x <= targetBounds.xUpper && currentPoint.y >= targetBounds.yLower {
		currentPoint.x += x
		currentPoint.y += y
		if x != 0 {
			x--
		}
		y--

		validX := currentPoint.x >= targetBounds.xLower && currentPoint.x <= targetBounds.xUpper
		validY := currentPoint.y >= targetBounds.yLower && currentPoint.y <= targetBounds.yUpper
		if validX && validY {
			return true
		}
	}

	return false
}

func findValidLaunchVelocities(targetBounds TargetBounds) []Point {
	maxY := 0 - targetBounds.yLower
	minX := findMinX(targetBounds.xLower)

	hits := make([]Point, 0)
	for x := minX; x <= targetBounds.xUpper; x++ {
		for y := targetBounds.yLower; y <= maxY; y++ {
			if trajectoryHitsTarget(x, y, targetBounds) {
				hits = append(hits, Point{x: x, y: y})
			}
		}
	}

	return hits
}

func main() {
	start := time.Now()
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	targetBounds := parseArea(scanner)

	hits := findValidLaunchVelocities(targetBounds)

	elapsed := time.Since(start)
	fmt.Println("Highest Y value sum:", calculateTriangle(-1-targetBounds.yLower))
	fmt.Println("Number of distinct starting velocities that will hit:", len(hits))
	log.Printf("Time taken: %s", elapsed)
}
