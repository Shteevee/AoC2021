package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func parseSyntaxLines(scanner *bufio.Scanner) []string {
	syntaxLines := make([]string, 0)
	for scanner.Scan() {
		syntaxLines = append(syntaxLines, scanner.Text())
	}

	return syntaxLines
}

func push(stack []rune, char rune) []rune {
	return append(stack, char)
}

func pop(stack []rune) (rune, []rune) {
	lastIndex := len(stack) - 1
	return stack[lastIndex], stack[:lastIndex]
}

func peek(stack []rune) rune {
	return stack[len(stack)-1]
}

func isEmpty(stack []rune) bool {
	return len(stack) == 0
}

func buildSyntaxPairs() map[rune]rune {
	syntaxPairs := make(map[rune]rune)
	syntaxPairs[')'] = '('
	syntaxPairs['}'] = '{'
	syntaxPairs['>'] = '<'
	syntaxPairs[']'] = '['

	return syntaxPairs
}

func isOpenningSyntax(char rune, syntaxPairs map[rune]rune) bool {
	for _, v := range syntaxPairs {
		if v == char {
			return true
		}
	}

	return false
}

func isLegalSyntax(syntaxLine string, syntaxPairs map[rune]rune) (bool, []rune) {
	stack := make([]rune, 0)

	for _, char := range syntaxLine {
		if isOpenningSyntax(char, syntaxPairs) {
			stack = push(stack, char)
		} else if syntaxPairs[char] == peek(stack) {
			_, stack = pop(stack)
		} else {
			return false, stack
		}
	}

	return true, stack
}

func calculateAutoCompleteStacks(syntaxLines []string) [][]rune {
	syntaxPairs := buildSyntaxPairs()
	autoCompleteStacks := make([][]rune, 0)

	for _, syntaxLine := range syntaxLines {
		isLegal, syntaxStack := isLegalSyntax(syntaxLine, syntaxPairs)
		if isLegal {
			autoCompleteStacks = append(autoCompleteStacks, syntaxStack)
		}
	}

	return autoCompleteStacks
}

func calculateAutoCompleteScore(autoCompleteStack []rune) int {
	var char rune
	score := 0

	for !isEmpty(autoCompleteStack) {
		score *= 5
		char, autoCompleteStack = pop(autoCompleteStack)
		if char == '(' {
			score += 1
		} else if char == '[' {
			score += 2
		} else if char == '{' {
			score += 3
		} else if char == '<' {
			score += 4
		}
	}

	return score
}

func calculateMiddleAutoCompleteScore(syntaxLines []string) int {
	autoCompleteStacks := calculateAutoCompleteStacks(syntaxLines)

	scores := make([]int, 0)
	for _, autoCompleteStack := range autoCompleteStacks {
		scores = append(scores, calculateAutoCompleteScore(autoCompleteStack))
	}

	sort.Ints(scores)

	return scores[len(scores)/2]
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
	syntaxLines := parseSyntaxLines(scanner)

	middleAutoCompleteScore := calculateMiddleAutoCompleteScore(syntaxLines)

	elapsed := time.Since(start)
	fmt.Println(middleAutoCompleteScore)
	log.Printf("Time taken: %s", elapsed)
}
