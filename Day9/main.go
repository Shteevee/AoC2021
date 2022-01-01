package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

type Point struct {
	x int
	y int
}

func parseHeightMap(scanner *bufio.Scanner) [][]int {
	heightMap := make([][]int, 0)
	for scanner.Scan() {
		heightRow := make([]int, 0)
		for _, char := range scanner.Text() {
			heightRow = append(heightRow, int(char-'0'))
		}
		heightMap = append(heightMap, heightRow)
	}

	return heightMap
}

func findLowPoints(heightMap [][]int) []Point {
	lowPoints := make([]Point, 0)

	for i, heightRow := range heightMap {
		for j, currentHeight := range heightRow {
			isLeftGreater := j == 0 || heightMap[i][j-1] > currentHeight
			isRightGreater := j == len(heightRow)-1 || heightMap[i][j+1] > currentHeight
			isUpGreater := i == 0 || heightMap[i-1][j] > currentHeight
			isDownGreater := i == len(heightMap)-1 || heightMap[i+1][j] > currentHeight

			if isLeftGreater && isRightGreater && isUpGreater && isDownGreater {
				lowPoints = append(lowPoints, Point{
					x: j,
					y: i,
				})
			}
		}
	}

	return lowPoints
}

func explored(point Point, exploredPoints []Point) bool {
	for _, exploredPoint := range exploredPoints {
		if point == exploredPoint {
			return true
		}
	}

	return false
}

// assumes a regular map (rectangular shape)
func exploreBasin(point Point, exploredPoints []Point, heightMap [][]int) (int, []Point) {
	exploredPoints = append(exploredPoints, point)
	exploredBasinSize := 1
	canMoveLeft := point.x > 0 && heightMap[point.y][point.x-1] < 9
	canMoveRight := point.x < len(heightMap[0])-1 && heightMap[point.y][point.x+1] < 9
	canMoveUp := point.y > 0 && heightMap[point.y-1][point.x] < 9
	canMoveDown := point.y < len(heightMap)-1 && heightMap[point.y+1][point.x] < 9

	if canMoveLeft {
		var discoveredBasinSize int
		left := Point{x: point.x - 1, y: point.y}
		if !explored(left, exploredPoints) {
			discoveredBasinSize, exploredPoints = exploreBasin(left, exploredPoints, heightMap)
		}
		exploredBasinSize += discoveredBasinSize
	}
	if canMoveRight {
		var discoveredBasinSize int
		right := Point{x: point.x + 1, y: point.y}
		if !explored(right, exploredPoints) {
			discoveredBasinSize, exploredPoints = exploreBasin(right, exploredPoints, heightMap)
		}
		exploredBasinSize += discoveredBasinSize
	}
	if canMoveUp {
		var discoveredBasinSize int
		up := Point{x: point.x, y: point.y - 1}
		if !explored(up, exploredPoints) {
			discoveredBasinSize, exploredPoints = exploreBasin(up, exploredPoints, heightMap)
		}
		exploredBasinSize += discoveredBasinSize
	}
	if canMoveDown {
		var discoveredBasinSize int
		down := Point{x: point.x, y: point.y + 1}
		if !explored(down, exploredPoints) {
			discoveredBasinSize, exploredPoints = exploreBasin(down, exploredPoints, heightMap)
		}
		exploredBasinSize += discoveredBasinSize
	}

	return exploredBasinSize, exploredPoints
}

func findBasinSizes(lowPoints []Point, heightMap [][]int) []int {
	exploredPoints := make([]Point, 0)
	basinSizes := make([]int, 0)

	for _, lowPoint := range lowPoints {
		if !explored(lowPoint, exploredPoints) {
			basinSize, expPoints := exploreBasin(lowPoint, make([]Point, 0), heightMap)
			exploredPoints = append(exploredPoints, expPoints...)
			basinSizes = append(basinSizes, basinSize)
		}
	}

	return basinSizes
}

func calculate3LargestBasinProduct(heightMap [][]int) int {
	lowPoints := findLowPoints(heightMap)
	basinSizes := findBasinSizes(lowPoints, heightMap)
	sort.Ints(basinSizes)

	basinProduct := 1
	for i := 1; i <= 3; i++ {
		basinProduct = basinProduct * basinSizes[len(basinSizes)-i]
	}

	return basinProduct
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
	heightMap := parseHeightMap(scanner)

	basinProduct := calculate3LargestBasinProduct(heightMap)

	elapsed := time.Since(start)
	fmt.Println(basinProduct)
	log.Printf("Time taken: %s", elapsed)
}
