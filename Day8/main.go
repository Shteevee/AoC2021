package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type WireSignal struct {
	signalPatterns []string
	output         []string
}

func parseWireSignals(scanner *bufio.Scanner) []WireSignal {
	wireSignals := make([]WireSignal, 0)
	for scanner.Scan() {
		splitSignalAndOutput := strings.Split(scanner.Text(), " | ")
		wireSignals = append(wireSignals, WireSignal{
			signalPatterns: strings.Split(splitSignalAndOutput[0], " "),
			output:         strings.Split(splitSignalAndOutput[1], " "),
		})
	}

	return wireSignals
}

func overlap(str1 string, str2 string) int {
	overlap := 0
	for _, char := range str1 {
		if strings.ContainsRune(str2, char) {
			overlap++
		}
	}

	return overlap
}

func intPow10(x int) int {
	return int(math.Pow10(x))
}

func unorderedMatch(str1 string, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}

	for _, char := range str1 {
		if !strings.ContainsRune(str2, char) {
			return false
		}
	}

	return true
}

// need to create an initial pattern map
// so that we can work out the more complex cases
func createInitialPatternMap(signalPatterns []string) ([]string, []string) {
	patternMap := make([]string, len(signalPatterns))
	remainingPatterns := make([]string, 0)
	for _, signalPattern := range signalPatterns {
		if len(signalPattern) == 2 {
			patternMap[1] = signalPattern
		} else if len(signalPattern) == 3 {
			patternMap[7] = signalPattern
		} else if len(signalPattern) == 4 {
			patternMap[4] = signalPattern
		} else if len(signalPattern) == 7 {
			patternMap[8] = signalPattern
		} else {
			remainingPatterns = append(remainingPatterns, signalPattern)
		}
	}

	return patternMap, remainingPatterns
}

func computePatternMap(signalPatterns []string) []string {
	patternMap, remainingPatterns := createInitialPatternMap(signalPatterns)

	for _, signalPattern := range remainingPatterns {
		if len(signalPattern) == 5 {
			if overlap(signalPattern, patternMap[4]) == 2 {
				patternMap[2] = signalPattern
			} else if overlap(signalPattern, patternMap[7]) == len(patternMap[7]) {
				patternMap[3] = signalPattern
			} else {
				patternMap[5] = signalPattern
			}
		} else if len(signalPattern) == 6 {
			if overlap(signalPattern, patternMap[7]) == 2 {
				patternMap[6] = signalPattern
			} else if overlap(signalPattern, patternMap[4]) == len(patternMap[4]) {
				patternMap[9] = signalPattern
			} else {
				patternMap[0] = signalPattern
			}
		}
	}

	return patternMap
}

func calculateWireOutput(wireSignal WireSignal) int {
	patternMap := computePatternMap(wireSignal.signalPatterns)

	output := 0
	for digitIndex, digit := range wireSignal.output {
		for patternIndex, pattern := range patternMap {
			if unorderedMatch(digit, pattern) {
				output = output + patternIndex*intPow10(len(wireSignal.output)-(digitIndex+1))
			}
		}
	}

	return output
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
	wireSignals := parseWireSignals(scanner)

	totalOutputs := 0
	for _, wireSignal := range wireSignals {
		totalOutputs = totalOutputs + calculateWireOutput(wireSignal)
	}

	fmt.Println("Sum of output digits: ", totalOutputs)
}
