package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Octopus struct {
	energy     int
	hasFlashed bool
}

type FlashPoint struct {
	x int
	y int
}

func parseOctopusGrid(scanner *bufio.Scanner) [][]Octopus {
	octopusGrid := make([][]Octopus, 0)
	for scanner.Scan() {
		octopusRow := make([]Octopus, 0)
		for _, octopus := range scanner.Text() {
			octopusRow = append(
				octopusRow,
				Octopus{
					energy:     int(octopus - '0'),
					hasFlashed: false,
				},
			)
		}
		octopusGrid = append(octopusGrid, octopusRow)
	}

	return octopusGrid
}

func increaseAllEnergyLevels(octopusGrid [][]Octopus) ([]FlashPoint, [][]Octopus) {
	flashed := make([]FlashPoint, 0)
	for i, octupusRow := range octopusGrid {
		for j := range octupusRow {
			if octopusGrid[i][j].energy == 9 {
				octopusGrid[i][j].energy = 0
				octopusGrid[i][j].hasFlashed = true
				flashed = append(flashed, FlashPoint{x: j, y: i})
			} else {
				octopusGrid[i][j].energy++
			}
		}
	}

	return flashed, octopusGrid
}

func incrementSurroundingEnergyLevels(point FlashPoint, octopusGrid [][]Octopus) ([]FlashPoint, [][]Octopus) {
	flashed := make([]FlashPoint, 0)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			y := point.y + i
			x := point.x + j
			validIndex := y >= 0 && y <= len(octopusGrid)-1 && x >= 0 && x <= len(octopusGrid[0])-1
			if validIndex && !octopusGrid[y][x].hasFlashed {
				//fmt.Println(x, y)
				if octopusGrid[y][x].energy == 9 {
					octopusGrid[y][x].energy = 0
					octopusGrid[y][x].hasFlashed = true
					flashed = append(flashed, FlashPoint{x: x, y: y})
				} else {
					octopusGrid[y][x].energy++
				}
			}
		}
	}

	return flashed, octopusGrid
}

func step(octopusGrid [][]Octopus) (int, [][]Octopus) {
	var flashPoints []FlashPoint
	flashPoints, octopusGrid = increaseAllEnergyLevels(octopusGrid)

	for i := 0; i < len(flashPoints); i++ {
		var newFlashPoints []FlashPoint
		newFlashPoints, octopusGrid = incrementSurroundingEnergyLevels(flashPoints[i], octopusGrid)
		flashPoints = append(flashPoints, newFlashPoints...)
	}

	return len(flashPoints), octopusGrid
}

func resetHasFlashed(octopusGrid [][]Octopus) [][]Octopus {
	for i, octopusRow := range octopusGrid {
		for j := range octopusRow {
			if octopusGrid[i][j].hasFlashed {
				octopusGrid[i][j].hasFlashed = false
			}
		}
	}

	return octopusGrid
}

func allHaveFlashed(octopusGrid [][]Octopus) bool {
	for i, octopusRow := range octopusGrid {
		for j := range octopusRow {
			if !octopusGrid[i][j].hasFlashed {
				return false
			}
		}
	}

	return true
}

func findFirstAllFlash(octopusGrid [][]Octopus) (int, int) {
	totalFlashes := 0
	stepCount := 0
	for !allHaveFlashed(octopusGrid) {
		stepCount++
		flashes := 0
		octopusGrid = resetHasFlashed(octopusGrid)
		flashes, octopusGrid = step(octopusGrid)
		totalFlashes += flashes
	}

	return stepCount, totalFlashes
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
	octopusGrid := parseOctopusGrid(scanner)

	steps, totalFlashes := findFirstAllFlash(octopusGrid)

	elapsed := time.Since(start)
	fmt.Println("All the octopuses flashed at step", steps)
	fmt.Println("No. of flashes:", totalFlashes)
	log.Printf("Time taken: %s", elapsed)
}
