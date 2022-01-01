package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type SnailfishNumber struct {
	Value int              `json:"Value"`
	Left  *SnailfishNumber `json:"Left"`
	Right *SnailfishNumber `json:"Right"`
	Depth int              `json:"Depth"`
}

func (sn *SnailfishNumber) String() string {
	if sn.Left == nil && sn.Right == nil {
		return fmt.Sprintf("%d", sn.Value)
	}
	return fmt.Sprintf("[%v,%v]", sn.Left, sn.Right)
}

func findSnailfishNumberSplitIndex(line string) int {
	bracketCount := 0
	for i, char := range line {
		if char == ',' && bracketCount == 0 {
			return i
		} else if char == '[' {
			bracketCount++
		} else if char == ']' {
			bracketCount--
		}
	}

	return -1
}

func parseSnailfishNumber(line string, Depth int) *SnailfishNumber {
	snailfishNumber := SnailfishNumber{}
	snailfishNumber.Depth = Depth

	if strings.Contains(line, ",") {
		line = strings.TrimPrefix(line, "[")
		line = strings.TrimSuffix(line, "]")
		splitIndex := findSnailfishNumberSplitIndex(line)

		snailfishNumber.Left = parseSnailfishNumber(line[:splitIndex], Depth+1)
		snailfishNumber.Right = parseSnailfishNumber(line[splitIndex+1:], Depth+1)
	} else {
		Value, _ := strconv.Atoi(line)
		snailfishNumber.Value = Value
	}

	return &snailfishNumber
}

func parseSnailfishNumbers(scanner *bufio.Scanner) []*SnailfishNumber {
	snailfishNumbers := make([]*SnailfishNumber, 0)
	for scanner.Scan() {
		snailfishNumbers = append(snailfishNumbers, parseSnailfishNumber(scanner.Text(), 0))
	}
	return snailfishNumbers
}

func postOrderTraverse(snailfishNumber *SnailfishNumber) []*SnailfishNumber {
	snailfishNumbers := make([]*SnailfishNumber, 0)
	if snailfishNumber.Left != nil && snailfishNumber.Right != nil {
		snailfishNumbers = append(snailfishNumbers, postOrderTraverse(snailfishNumber.Left)...)
		snailfishNumbers = append(snailfishNumbers, postOrderTraverse(snailfishNumber.Right)...)
	}
	snailfishNumbers = append(snailfishNumbers, snailfishNumber)

	return snailfishNumbers
}

func incrementDepth(sn *SnailfishNumber) {
	sn.Depth++
	if sn.Left != nil && sn.Right != nil {
		incrementDepth(sn.Left)
		incrementDepth(sn.Right)
	}
}

func add(left *SnailfishNumber, right *SnailfishNumber) *SnailfishNumber {
	leftCpy := copyNumber(left)
	rightCpy := copyNumber(right)
	incrementDepth(&leftCpy)
	incrementDepth(&rightCpy)
	return &SnailfishNumber{
		Left:  &leftCpy,
		Right: &rightCpy,
	}
}

func canExplode(snailfishNumber *SnailfishNumber) bool {
	explode := false
	postOrderNumbers := postOrderTraverse(snailfishNumber)
	for _, number := range postOrderNumbers {
		explode = explode || number.Depth == 4 && number.Left != nil && number.Right != nil
	}

	return explode
}

func canSplit(snailfishNumber *SnailfishNumber) bool {
	if snailfishNumber.Left != nil && snailfishNumber.Right != nil {
		return canSplit(snailfishNumber.Left) || canSplit(snailfishNumber.Right)
	}

	return snailfishNumber.Value >= 10
}

func split(snailfishNumber *SnailfishNumber) {
	postOrderNumbers := postOrderTraverse(snailfishNumber)

	for i := range postOrderNumbers {
		if postOrderNumbers[i].Value >= 10 {
			half := postOrderNumbers[i].Value / 2
			postOrderNumbers[i].Left = &SnailfishNumber{
				Value: half,
				Depth: postOrderNumbers[i].Depth + 1,
			}
			postOrderNumbers[i].Right = &SnailfishNumber{
				Value: half + postOrderNumbers[i].Value%2,
				Depth: postOrderNumbers[i].Depth + 1,
			}
			postOrderNumbers[i].Value = 0
			return
		}
	}
}

func findExplodeIndex(snailfishNumbers []*SnailfishNumber) int {
	for i, snailfishNumber := range snailfishNumbers {
		if snailfishNumber.Depth == 4 && snailfishNumber.Left != nil && snailfishNumber.Right != nil {
			return i
		}
	}

	return -1
}

func explodeToLeft(postOrderNumbers []*SnailfishNumber, explodeIndex int) {
	for i := explodeIndex - 3; i >= 0; i-- {
		if postOrderNumbers[i].Left == nil && postOrderNumbers[i].Right == nil {
			postOrderNumbers[i].Value += postOrderNumbers[explodeIndex].Left.Value
			return
		}
	}
}

func explodeToRight(postOrderNumbers []*SnailfishNumber, explodeIndex int) {
	for i := explodeIndex + 1; i < len(postOrderNumbers); i++ {
		if postOrderNumbers[i].Left == nil && postOrderNumbers[i].Right == nil {
			postOrderNumbers[i].Value += postOrderNumbers[explodeIndex].Right.Value
			return
		}
	}
}

func explode(snailfishNumber *SnailfishNumber) {
	postOrderNumbers := postOrderTraverse(snailfishNumber)
	explodeIndex := findExplodeIndex(postOrderNumbers)
	explodeToLeft(postOrderNumbers, explodeIndex)
	explodeToRight(postOrderNumbers, explodeIndex)

	postOrderNumbers[explodeIndex].Value = 0
	postOrderNumbers[explodeIndex].Left = nil
	postOrderNumbers[explodeIndex].Right = nil
}

func reduce(snailfishNumber *SnailfishNumber) *SnailfishNumber {
	sn := copyNumber(snailfishNumber)
	shouldExplode := canExplode(&sn)
	shouldSplit := canSplit(&sn)

	for shouldExplode || shouldSplit {
		if shouldExplode {
			explode(&sn)
		} else if shouldSplit {
			split(&sn)
		}

		shouldExplode = canExplode(&sn)
		shouldSplit = canSplit(&sn)
	}

	return &sn
}

func doHomework(snailfishNumbers []*SnailfishNumber) *SnailfishNumber {
	total := snailfishNumbers[0]
	for _, snailfishNumber := range snailfishNumbers[1:] {
		total = add(total, snailfishNumber)
		total = reduce(total)
	}

	return total
}

func calculateMagnitude(sn *SnailfishNumber) int {
	total := 0
	if sn.Left != nil && sn.Right != nil {
		total += 3 * calculateMagnitude(sn.Left)
		total += 2 * calculateMagnitude(sn.Right)
	} else {
		total += sn.Value
	}

	return total
}

// should really have made pure functions but this was quicker
func copyNumber(sn *SnailfishNumber) SnailfishNumber {
	data, err := json.Marshal(sn)
	if err != nil {
		panic(err)
	}

	var newSn SnailfishNumber
	if err := json.Unmarshal(data, &newSn); err != nil {
		panic(err)
	}

	return newSn
}

func findLargestMagnitude(snailfishNumbers []*SnailfishNumber) int {
	max := math.MinInt
	for i, x := range snailfishNumbers {
		for j, y := range snailfishNumbers {
			if i != j {
				result := add(x, y)
				result = reduce(result)
				mag := calculateMagnitude(result)
				if mag > max {
					max = mag
				}
			}
		}
	}

	return max
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
	snailfishNumbers := parseSnailfishNumbers(scanner)

	total := doHomework(snailfishNumbers)
	magnitude := calculateMagnitude(total)
	largestMagnitude := findLargestMagnitude(snailfishNumbers)

	elapsed := time.Since(start)
	fmt.Println(magnitude)
	fmt.Println(largestMagnitude)
	log.Printf("Time taken: %s", elapsed)
}
