package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseCrabs(scanner *bufio.Scanner) []int {
	scanner.Scan()
	stringCrabs := strings.Split(scanner.Text(), ",")
	crabs := make([]int, len(stringCrabs))
	for i, strCrab := range stringCrabs {
		crab, _ := strconv.Atoi(strCrab)
		crabs[i] = crab
	}

	return crabs
}

func max(arr []int) int {
	max := arr[0]
	for _, val := range arr[1:] {
		if val > max {
			max = val
		}
	}

	return max
}

func min(arr []int) int {
	min := arr[0]
	for _, val := range arr[1:] {
		if val < min {
			min = val
		}
	}

	return min
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func calculateTriangle(x int) int {
	return (x * (x + 1)) / 2
}

func findShortestTotalDistance(crabs []int) int {
	maxCrab := max(crabs)
	minCrab := min(crabs)

	cheapestFuel := math.MaxInt32
	for position := minCrab; position < maxCrab; position++ {
		fuel := 0
		for _, crab := range crabs {
			fuel = fuel + calculateTriangle(abs(crab-position))
		}
		if fuel < cheapestFuel {
			cheapestFuel = fuel
		}
	}

	return cheapestFuel
}

func main() {
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
	crabs := parseCrabs(scanner)

	cheapestFuel := findShortestTotalDistance(crabs)

	fmt.Println("Cheapest fuel: ", cheapestFuel)
}
