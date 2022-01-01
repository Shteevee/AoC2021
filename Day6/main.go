package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//the word fish being it's own plural is very annoying
func parseFish(scanner *bufio.Scanner) []int {
	scanner.Scan()
	stringFish := strings.Split(scanner.Text(), ",")
	fish := make([]int, 0)
	for _, strFish := range stringFish {
		newFish, _ := strconv.Atoi(strFish)
		fish = append(fish, newFish)
	}

	return fish
}

func populateFishMap(fish []int, fishMap []int) []int {
	for _, aFish := range fish {
		fishMap[aFish]++
	}
	return fishMap
}

//uses a "fishMap" which is essentially just a bucket
//by fish timer value
func passDays(fish []int, days int) []int {
	fishMap := make([]int, 9)
	fishMap = populateFishMap(fish, fishMap)
	for i := 0; i < days; i++ {
		temp := fishMap[0]
		fishMap = fishMap[1:]
		fishMap = append(fishMap, temp)
		fishMap[6] = fishMap[6] + temp
	}

	return fishMap
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

	fish := parseFish(scanner)

	fishMap := passDays(fish, 256)

	totalFish := 0
	for _, numberOfFish := range fishMap {
		totalFish = totalFish + numberOfFish
	}

	fmt.Println("Number of fish: ", totalFish)
}
