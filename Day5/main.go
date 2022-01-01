package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Line struct {
	start Point
	end   Point
}

func parsePoint(strPoint string) Point {
	splitPoint := strings.Split(strPoint, ",")
	x, _ := strconv.Atoi(splitPoint[0])
	y, _ := strconv.Atoi(splitPoint[1])
	point := Point{
		x: x,
		y: y,
	}
	return point
}

func parseLines(scanner *bufio.Scanner) []Line {
	lines := make([]Line, 0)
	for scanner.Scan() {
		splitText := strings.Split(scanner.Text(), " -> ")
		line := Line{
			start: parsePoint(splitText[0]),
			end:   parsePoint(splitText[1]),
		}
		lines = append(lines, line)
	}
	return lines
}

func filterHorizontalAndVerticalLines(lines []Line) []Line {
	n := 0
	for _, line := range lines {
		if line.start.x == line.end.x || line.start.y == line.end.y {
			lines[n] = line
			n++
		}
	}
	return lines[:n]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

//this will only work for 45 degree lines
func spreadLine(line Line) []Point {
	vectorX := line.end.x - line.start.x
	vectorY := line.end.y - line.start.y

	xs := make([]int, 0)
	ys := make([]int, 0)
	if vectorX > 0 {
		for i := 0; i <= vectorX; i++ {
			xs = append(xs, line.start.x+i)
		}
	} else if vectorX < 0 {
		for i := 0; i >= vectorX; i-- {
			xs = append(xs, line.start.x+i)
		}
	} else {
		for i := 0; i <= abs(vectorY); i++ {
			xs = append(xs, line.start.x)
		}
	}

	if vectorY > 0 {
		for i := 0; i <= vectorY; i++ {
			ys = append(ys, line.start.y+i)
		}
	} else if vectorY < 0 {
		for i := 0; i >= vectorY; i-- {
			ys = append(ys, line.start.y+i)
		}
	} else {
		for i := 0; i <= abs(vectorX); i++ {
			ys = append(ys, line.start.y)
		}
	}

	points := make([]Point, 0)
	for i := 0; i < len(xs); i++ {
		point := Point{
			x: xs[i],
			y: ys[i],
		}
		points = append(points, point)
	}

	return points
}

func createPointsMap(lines []Line) map[Point]int {
	pointsMap := make(map[Point]int)
	for _, line := range lines {
		points := spreadLine(line)
		for _, point := range points {
			pointsMap[point]++
		}
	}

	return pointsMap
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

	lines := parseLines(scanner)

	//lines = filterHorizontalAndVerticalLines(lines)

	pointsMap := createPointsMap(lines)

	overlapPoints := 0
	for _, v := range pointsMap {
		if v > 1 {
			overlapPoints++
		}
	}

	fmt.Println("Number of overlaps: ", overlapPoints)
}
