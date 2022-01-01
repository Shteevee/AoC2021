package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

func parseTemplateAndPairs(scanner *bufio.Scanner) (map[string]int, map[string]string) {
	templateMap := make(map[string]int, 0)
	pairMap := make(map[string]string)
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "->") && scanner.Text() != "" {
			for i := 0; i < len(scanner.Text())-1; i++ {
				templateMap[string(scanner.Text()[i])+string(scanner.Text()[i+1])]++
			}
		} else if strings.Contains(scanner.Text(), "->") {
			splitPairs := strings.Split(scanner.Text(), " -> ")
			pairMap[splitPairs[0]] = splitPairs[1]
		}
	}

	return templateMap, pairMap
}

func pairInsertion(polyTemplate map[string]int, pairMap map[string]string) map[string]int {
	newPolyTemplate := make(map[string]int)

	for pair, count := range polyTemplate {
		newPolyTemplate[string(pair[0])+pairMap[pair]] += count
		newPolyTemplate[pairMap[pair]+string(pair[1])] += count
	}

	return newPolyTemplate
}

func performPairInsertions(polyTemplateMap map[string]int, pairMap map[string]string, steps int) map[string]int {
	for i := 0; i < steps; i++ {
		polyTemplateMap = pairInsertion(polyTemplateMap, pairMap)
	}

	return polyTemplateMap
}

func countElements(polymerMap map[string]int) map[rune]int {
	elementCounts := make(map[rune]int)
	for pair, count := range polymerMap {
		for _, char := range pair {
			elementCounts[char] += count
		}
	}

	for element, count := range elementCounts {
		elementCounts[element] = count/2 + count%2
	}

	return elementCounts
}

func differenceMaxAndMinOccurences(elementCounts map[rune]int) int {
	min := math.MaxInt
	max := math.MinInt
	for _, count := range elementCounts {
		if count > max {
			max = count
		} else if count < min {
			min = count
		}
	}

	return max - min
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
	polyTemplateMap, pairMap := parseTemplateAndPairs(scanner)

	polyTemplateMap = performPairInsertions(polyTemplateMap, pairMap, 40)
	elementCounts := countElements(polyTemplateMap)
	diffMaxMin := differenceMaxAndMinOccurences(elementCounts)

	total := 1
	for _, count := range polyTemplateMap {
		total += count
	}

	elapsed := time.Since(start)
	fmt.Println(elementCounts)
	fmt.Println("Difference between max and min occurences: ", diffMaxMin)
	log.Printf("Time taken: %s", elapsed)
}
