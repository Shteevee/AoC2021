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

type Fold struct {
	axis  string
	index int
}

type Point struct {
	x int
	y int
}

func parsePaperAndInstructions(scanner *bufio.Scanner) (map[Point]bool, []Fold) {
	points := make(map[Point]bool, 0)
	folds := make([]Fold, 0)

	for scanner.Scan() {
		if strings.ContainsRune(scanner.Text(), ',') {
			strPoints := strings.Split(scanner.Text(), ",")
			x, _ := strconv.Atoi(strPoints[0])
			y, _ := strconv.Atoi(strPoints[1])
			newPoint := Point{
				x: x,
				y: y,
			}
			points[newPoint] = true
		} else if strings.HasPrefix(scanner.Text(), "fold along") {
			instr := strings.TrimPrefix(scanner.Text(), "fold along ")
			parts := strings.Split(instr, "=")
			index, _ := strconv.Atoi(parts[1])
			folds = append(
				folds,
				Fold{
					axis:  parts[0],
					index: index,
				},
			)
		}
	}

	return points, folds
}

func findMaxY(points map[Point]bool) int {
	max := 0
	for point := range points {
		if point.y > max {
			max = point.y
		}
	}

	return max
}

func findMaxX(points map[Point]bool) int {
	max := 0
	for point := range points {
		if point.x > max {
			max = point.x
		}
	}

	return max
}

func foldX(point Point, index int, max int) Point {
	if point.x > index {
		return Point{x: max - point.x, y: point.y}
	}

	return point
}

func foldY(point Point, index int, max int) Point {
	if point.y > index {
		return Point{x: point.x, y: max - point.y}
	}

	return point
}

func performFold(points map[Point]bool, fold Fold) map[Point]bool {
	newPoints := make(map[Point]bool)
	if fold.axis == "x" {
		maxX := findMaxX(points)
		for point := range points {
			foldPoint := foldX(point, fold.index, maxX)
			if foldPoint.x >= 0 && foldPoint.x < fold.index {
				newPoints[foldPoint] = true
			}
		}
	} else {
		maxY := findMaxY(points)
		for point := range points {
			foldPoint := foldY(point, fold.index, maxY)
			if foldPoint.y >= 0 && foldPoint.y < fold.index {
				newPoints[foldPoint] = true
			}
		}
	}

	return newPoints
}

func performFolds(points map[Point]bool, folds []Fold) map[Point]bool {
	for _, fold := range folds {
		points = performFold(points, fold)
		fmt.Println(fold, "No. of points: ", len(points))
	}

	return points
}

func displayPoints(points map[Point]bool) {
	maxX := findMaxX(points) + 1
	maxY := findMaxY(points) + 1

	paper := make([][]string, maxY)
	for i := range paper {
		paper[i] = make([]string, maxX)
		for str := range paper[i] {
			paper[i][str] = "."
		}
	}

	for point := range points {
		paper[point.y][point.x] = "#"
	}

	for _, line := range paper {
		fmt.Println(line)
	}
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
	points, folds := parsePaperAndInstructions(scanner)

	points = performFolds(points, folds)

	displayPoints(points)

	elapsed := time.Since(start)
	fmt.Println("Final number of points: ", len(points))
	log.Printf("Time taken: %s", elapsed)
}
