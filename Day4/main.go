package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseNumbers(line string, splitter string) []int {
	numbersSplit := strings.Split(line, splitter)
	var numbers []int
	for _, drawnNumber := range numbersSplit {
		if drawnNumber != "" {
			parsedNumber, _ := strconv.Atoi(drawnNumber)
			numbers = append(numbers, parsedNumber)
		}
	}

	return numbers
}

func createBoardArray(scanner *bufio.Scanner) [][][]int {
	var boardArray [][][]int
	var currentBoard [][]int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			boardArray = append(boardArray, currentBoard)
			currentBoard = *new([][]int)
		} else {
			currentBoard = append(currentBoard, parseNumbers(line, " "))
		}
	}
	boardArray = append(boardArray, currentBoard)

	return boardArray
}

//Assumes square, regular inputs
func createMatchTracker(boardArray [][][]int) [][][]bool {
	var matchTracker [][][]bool
	for i := 0; i < len(boardArray); i++ {
		matchBoard := make([][]bool, len(boardArray[0]))
		for j := range matchBoard {
			matchBoard[j] = make([]bool, len(boardArray[0]))
		}
		matchTracker = append(matchTracker, matchBoard)
	}

	return matchTracker
}

func markDrawnNumber(boardArray [][][]int, matchTracker [][][]bool, drawnNumber int) ([][][]int, [][][]bool) {
	for i, board := range boardArray {
		for j, row := range board {
			for k, elem := range row {
				if elem == drawnNumber {
					matchTracker[i][j][k] = true
				}
			}
		}
	}

	return boardArray, matchTracker
}

func findWinner(matchTracker [][][]bool) int {
	for i, matchBoard := range matchTracker {
		if checkBingo(matchBoard) {
			return i
		}
	}

	return -1
}

func contains(winners []int, possibleWinner int) bool {
	for _, winner := range winners {
		if winner == possibleWinner {
			return true
		}
	}

	return false
}

func findNewWinners(matchTracker [][][]bool, winners []int) []int {
	var newWinners []int
	for i, matchBoard := range matchTracker {
		if checkBingo(matchBoard) && !contains(winners, i) {
			newWinners = append(newWinners, i)
		}
	}

	return newWinners
}

func checkBingo(matchBoard [][]bool) bool {
	for _, row := range matchBoard {
		bingo := true
		for _, elem := range row {
			bingo = bingo && elem
		}
		if bingo {
			return true
		}
	}

	for i := range matchBoard {
		bingo := true
		for j := range matchBoard {
			bingo = bingo && matchBoard[j][i]
		}
		if bingo {
			return true
		}
	}

	return false
}

func playBingoToLose(drawnNumbers []int, boardArray [][][]int, matchTracker [][][]bool) (int, [][]int, [][]bool) {
	winners := make([]int, 0, len(boardArray))
	var lastCalled int
	var loserIndex int
	for _, drawnNumber := range drawnNumbers {
		boardArray, matchTracker = markDrawnNumber(boardArray, matchTracker, drawnNumber)

		newWinners := findNewWinners(matchTracker, winners)
		if len(newWinners) > 0 {
			fmt.Println("Winners! ", newWinners)
			winners = append(winners, newWinners...)
			if len(winners) == len(boardArray) {
				lastCalled = drawnNumber
				loserIndex = newWinners[0]
				fmt.Println("Losers: ", newWinners)
				break
			}
		}

		fmt.Println(winners)
	}

	return lastCalled, boardArray[loserIndex], matchTracker[loserIndex]
}

func playBingo(drawnNumbers []int, boardArray [][][]int, matchTracker [][][]bool) (int, [][]int, [][]bool) {
	winnerIndex := -1
	var lastCalled int
	for _, drawnNumber := range drawnNumbers {
		boardArray, matchTracker = markDrawnNumber(boardArray, matchTracker, drawnNumber)

		winnerIndex = findWinner(matchTracker)
		if winnerIndex > -1 {
			lastCalled = drawnNumber
			break
		}
	}

	return lastCalled, boardArray[winnerIndex], matchTracker[winnerIndex]
}

func calculateScore(lastCalled int, board [][]int, matches [][]bool) int {
	total := 0
	for i, row := range matches {
		for j, elem := range row {
			if !elem {
				total = total + board[i][j]
			}
		}
	}

	return total * lastCalled
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

	scanner.Scan()
	drawnNumbers := parseNumbers(scanner.Text(), ",")

	scanner.Scan()
	boardArray := createBoardArray(scanner)

	matchTracker := createMatchTracker(boardArray)

	lastCalled, winningBoard, winningMatch := playBingoToLose(drawnNumbers, boardArray, matchTracker)

	score := calculateScore(lastCalled, winningBoard, winningMatch)

	fmt.Println("Last call: ", lastCalled)
	fmt.Println("Winning board: ", winningBoard)
	fmt.Println("Match board: ", winningMatch)
	fmt.Println("Drawn numbers: ", drawnNumbers)
	fmt.Println("Score: ", score)
}
