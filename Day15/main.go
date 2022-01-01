package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type QueueItem struct {
	value    Point
	priority int
	index    int
}

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*QueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func push(pq *PriorityQueue, x int, y int, distance int) {
	heap.Push(
		pq,
		&QueueItem{
			value:    Point{x: x, y: y},
			priority: distance,
		},
	)
}

type Point struct {
	x int
	y int
}

type RiskNode struct {
	risk     int
	distance int
	visited  bool
}

func parseRiskMap(scanner *bufio.Scanner) [][]RiskNode {
	riskMap := make([][]RiskNode, 0)
	for scanner.Scan() {
		riskRow := make([]RiskNode, 0)
		for _, char := range scanner.Text() {
			riskRow = append(
				riskRow,
				RiskNode{
					risk:     int(char - '0'),
					distance: math.MaxInt,
					visited:  false,
				},
			)
		}
		riskMap = append(riskMap, riskRow)
	}

	return riskMap
}

func incrementRiskMap(riskMap [][]RiskNode) [][]RiskNode {
	for i, riskRow := range riskMap {
		for j, riskNode := range riskRow {
			if riskNode.risk == 9 {
				riskMap[i][j].risk = 1
			} else {
				riskMap[i][j].risk++
			}
		}
	}

	return riskMap
}

func copyRiskMap(riskMap [][]RiskNode) [][]RiskNode {
	duplicate := make([][]RiskNode, len(riskMap))
	for i := range riskMap {
		duplicate[i] = make([]RiskNode, len(riskMap[i]))
		copy(duplicate[i], riskMap[i])
	}

	return duplicate
}

func expandCaveSystem(riskMap [][]RiskNode) [][]RiskNode {
	incrementedRiskMap := copyRiskMap(riskMap)
	for k := 0; k < 4; k++ {
		incrementedRiskMap = incrementRiskMap(incrementedRiskMap)
		for i := range riskMap {
			riskMap[i] = append(riskMap[i], incrementedRiskMap[i]...)
		}
	}

	incrementedRiskMap = copyRiskMap(riskMap)
	for k := 0; k < 4; k++ {
		incrementedRiskMap = incrementRiskMap(copyRiskMap(incrementedRiskMap))
		riskMap = append(riskMap, incrementedRiskMap...)
	}

	return riskMap
}

func findShortestPath(riskMap [][]RiskNode) [][]RiskNode {
	exploring := true
	x := 0
	y := 0
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for exploring {
		shouldExploreLeft := x > 0 && !riskMap[y][x-1].visited && riskMap[y][x-1].distance > riskMap[y][x-1].risk+riskMap[y][x].distance
		shouldExploreRight := x < len(riskMap[0])-1 && !riskMap[y][x+1].visited && riskMap[y][x+1].distance > riskMap[y][x+1].risk+riskMap[y][x].distance
		shouldExploreUp := y > 0 && !riskMap[y-1][x].visited && riskMap[y-1][x].distance > riskMap[y-1][x].risk+riskMap[y][x].distance
		shouldExploreDown := y < len(riskMap)-1 && !riskMap[y+1][x].visited && riskMap[y+1][x].distance > riskMap[y+1][x].risk+riskMap[y][x].distance

		if shouldExploreLeft {
			riskMap[y][x-1].distance = riskMap[y][x-1].risk + riskMap[y][x].distance
			push(&pq, x-1, y, riskMap[y][x-1].distance)
		}
		if shouldExploreRight {
			riskMap[y][x+1].distance = riskMap[y][x+1].risk + riskMap[y][x].distance
			push(&pq, x+1, y, riskMap[y][x+1].distance)
		}
		if shouldExploreUp {
			riskMap[y-1][x].distance = riskMap[y-1][x].risk + riskMap[y][x].distance
			push(&pq, x, y-1, riskMap[y-1][x].distance)
		}
		if shouldExploreDown {
			riskMap[y+1][x].distance = riskMap[y+1][x].risk + riskMap[y][x].distance
			push(&pq, x, y+1, riskMap[y+1][x].distance)
		}
		riskMap[y][x].visited = true
		nextPoint := heap.Pop(&pq).(*QueueItem)
		x = nextPoint.value.x
		y = nextPoint.value.y
		if x == len(riskMap[0])-1 && y == len(riskMap)-1 {
			exploring = false
		}
	}

	return riskMap
}

func findShortestPathLength(riskMap [][]RiskNode) int {
	riskMap = findShortestPath(riskMap)

	return riskMap[len(riskMap)-1][len(riskMap[0])-1].distance
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
	riskMap := parseRiskMap(scanner)
	riskMap = expandCaveSystem(riskMap)

	riskMap[0][0].distance = 0
	riskMap[0][0].visited = true
	pathLength := findShortestPathLength(riskMap)

	elapsed := time.Since(start)
	fmt.Println("Shortest path length:", pathLength)
	log.Printf("Time taken: %s", elapsed)
}
