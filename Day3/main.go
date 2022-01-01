package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

//It's pretty gross down there
//Heed my warning

const bitsLength int = 12

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func lineToBits(line string) [bitsLength]int {
	var number [bitsLength]int
	for i := 0; i < len(line); i++ {
		number[i], _ = strconv.Atoi(string(line[i]))
	}

	return number
}

func filterByBit(bitsArray [][bitsLength]int, index int, filterValue int) [][bitsLength]int {
	n := 0
	for _, bits := range bitsArray {
		if bits[index] == filterValue {
			bitsArray[n] = bits
			n++
		}
	}

	if len(bitsArray[:n]) == 0 {
		return bitsArray
	}

	return bitsArray[:n]
}

func incrementCounters(bits [bitsLength]int, index int, zeros int, ones int) (int, int) {
	if bits[index] == 1 {
		ones++
	} else if bits[index] == 0 {
		zeros++
	}

	return zeros, ones
}

func binToDec(bits []int) int {
	total := 0
	for i := 0; i < bitsLength; i++ {
		if bits[i] == 1 {
			total = total + powInt(2, bitsLength-1-i)
		}
	}

	return total
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

	var bitsArray [][bitsLength]int

	for scanner.Scan() {
		line := scanner.Text()
		bitsArray = append(bitsArray, lineToBits(line))
	}

	oxygenGenCandidates := make([][bitsLength]int, len(bitsArray))
	copy(oxygenGenCandidates, bitsArray)
	co2ScrubCandidates := make([][bitsLength]int, len(bitsArray))
	copy(co2ScrubCandidates, bitsArray)

	for i := 0; i < bitsLength; i++ {
		oxygenZeros := 0
		oxygenOnes := 0
		co2Zeros := 0
		co2Ones := 0
		for j := 0; j < len(oxygenGenCandidates); j++ {
			oxygenZeros, oxygenOnes = incrementCounters(oxygenGenCandidates[j], i, oxygenZeros, oxygenOnes)
		}
		for j := 0; j < len(co2ScrubCandidates); j++ {
			co2Zeros, co2Ones = incrementCounters(co2ScrubCandidates[j], i, co2Zeros, co2Ones)
		}

		fmt.Println("Oxygen:   0: ", oxygenZeros, "1: ", oxygenOnes)
		fmt.Println("Carbon:   0: ", co2Zeros, "1: ", co2Ones)
		fmt.Println()
		oxygenFilterValue := 1
		if oxygenZeros > oxygenOnes {
			oxygenFilterValue = 0
		}

		co2FilterValue := 0
		if co2Ones < co2Zeros {
			co2FilterValue = 1
		}

		oxygenGenCandidates = filterByBit(oxygenGenCandidates, i, oxygenFilterValue)
		co2ScrubCandidates = filterByBit(co2ScrubCandidates, i, co2FilterValue)
	}

	fmt.Println("Oxygen generator rating: ", oxygenGenCandidates, binToDec(oxygenGenCandidates[0][:]))
	fmt.Println("CO2 scrubber rating: ", co2ScrubCandidates, binToDec(co2ScrubCandidates[0][:]))
	fmt.Println("Result: ", binToDec(co2ScrubCandidates[0][:])*binToDec(oxygenGenCandidates[0][:]))
}
