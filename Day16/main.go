package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type Packet struct {
	version      int
	typeId       int
	lengthTypeId int
	value        int
	packets      []Packet
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func getHexBinMap() map[rune][]int {
	hexBinMap := make(map[rune][]int)
	hexBinMap['0'] = []int{0, 0, 0, 0}
	hexBinMap['1'] = []int{0, 0, 0, 1}
	hexBinMap['2'] = []int{0, 0, 1, 0}
	hexBinMap['3'] = []int{0, 0, 1, 1}
	hexBinMap['4'] = []int{0, 1, 0, 0}
	hexBinMap['5'] = []int{0, 1, 0, 1}
	hexBinMap['6'] = []int{0, 1, 1, 0}
	hexBinMap['7'] = []int{0, 1, 1, 1}
	hexBinMap['8'] = []int{1, 0, 0, 0}
	hexBinMap['9'] = []int{1, 0, 0, 1}
	hexBinMap['A'] = []int{1, 0, 1, 0}
	hexBinMap['B'] = []int{1, 0, 1, 1}
	hexBinMap['C'] = []int{1, 1, 0, 0}
	hexBinMap['D'] = []int{1, 1, 0, 1}
	hexBinMap['E'] = []int{1, 1, 1, 0}
	hexBinMap['F'] = []int{1, 1, 1, 1}

	return hexBinMap
}

func parseHex(scanner *bufio.Scanner) []int {
	hexBinMap := getHexBinMap()
	parsedHex := make([]int, 0)

	scanner.Scan()
	for _, char := range scanner.Text() {
		parsedHex = append(parsedHex, hexBinMap[char]...)
	}

	return parsedHex
}

func binToDec(binary []int) int {
	total := 0
	for i := 0; i < len(binary); i++ {
		if binary[i] == 1 {
			total = total + powInt(2, len(binary)-1-i)
		}
	}

	return total
}

func filterEndIndicators(binary []int) []int {
	newArray := make([]int, 0)
	for i, x := range binary {
		if i%5 != 0 {
			newArray = append(newArray, x)
		}
	}
	return newArray
}

func parsePacket(binary []int, i int) (Packet, int) {
	packet := Packet{}
	packet.version = binToDec(binary[i : i+3])
	packet.typeId = binToDec(binary[i+3 : i+6])

	if packet.typeId == 4 {
		start := i + 6
		i += 6 + 5
		for binary[i-5] != 0 {
			i += 5
		}
		packet.value = binToDec(filterEndIndicators(binary[start:i]))
		return packet, i
	} else {
		packet.lengthTypeId = binToDec(binary[i+6 : i+7])
		if packet.lengthTypeId == 0 {
			length := binToDec(binary[i+7 : i+22])
			i += 22
			start := i
			for i-start < length {
				newPacket, newIndex := parsePacket(binary, i)
				packet.packets = append(packet.packets, newPacket)
				i = newIndex
			}
			return packet, i
		} else {
			subPacketCount := binToDec(binary[i+7 : i+18])
			i += 18
			for len(packet.packets) < subPacketCount {
				newPacket, newIndex := parsePacket(binary, i)
				packet.packets = append(packet.packets, newPacket)
				i = newIndex
			}
			return packet, i
		}
	}
}

func sumVersions(packet Packet) int {
	total := 0
	if len(packet.packets) > 0 {
		for _, p := range packet.packets {
			total += sumVersions(p)
		}
	}

	return total + packet.version
}

func sum(packets []Packet) int {
	total := 0
	for _, p := range packets {
		total += evaluate(p)
	}
	return total
}

func product(packets []Packet) int {
	total := 1
	for _, p := range packets {
		total *= evaluate(p)
	}
	return total
}

func min(packets []Packet) int {
	min := math.MaxInt
	for _, p := range packets {
		packetValue := evaluate(p)
		if packetValue < min {
			min = packetValue
		}
	}
	return min
}

func max(packets []Packet) int {
	max := math.MinInt
	for _, p := range packets {
		packetValue := evaluate(p)
		if packetValue > max {
			max = packetValue
		}
	}
	return max
}

func greaterThan(packets []Packet) int {
	if evaluate(packets[0]) > evaluate(packets[1]) {
		return 1
	}

	return 0
}

func lessThan(packets []Packet) int {
	if evaluate(packets[0]) < evaluate(packets[1]) {
		return 1
	}

	return 0
}

func equal(packets []Packet) int {
	if evaluate(packets[0]) == evaluate(packets[1]) {
		return 1
	}

	return 0
}

func evaluate(packet Packet) int {
	if len(packet.packets) > 0 {
		switch packet.typeId {
		case 0:
			return sum(packet.packets)
		case 1:
			return product(packet.packets)
		case 2:
			return min(packet.packets)
		case 3:
			return max(packet.packets)
		case 5:
			return greaterThan(packet.packets)
		case 6:
			return lessThan(packet.packets)
		case 7:
			return equal(packet.packets)
		}
	}

	return packet.value
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
	binary := parseHex(scanner)

	packet, _ := parsePacket(binary, 0)
	versionSum := sumVersions(packet)
	packetValue := evaluate(packet)

	elapsed := time.Since(start)
	fmt.Println("Version sum:", versionSum)
	fmt.Println("Packet value:", packetValue)
	log.Printf("Time taken: %s", elapsed)
}
