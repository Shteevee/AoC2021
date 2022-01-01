package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func sum(nums [3]int) int {
	var total int
	for _, x := range nums {
		total = total + x
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

	counter := 0
	var prev [3]int
	var curr [3]int

	scanner.Scan()
	prev[0], _ = strconv.Atoi(scanner.Text())
	for i := 1; i < 3; i++ {
		scanner.Scan()
		prev[i], _ = strconv.Atoi(scanner.Text())
		curr[i-1], _ = strconv.Atoi(scanner.Text())
	}

	for scanner.Scan() {
		curr[2], _ = strconv.Atoi(scanner.Text())
		if sum(prev) < sum(curr) {
			counter++
		}
		prev = curr
		curr[0] = curr[1]
		curr[1] = curr[2]
	}

	fmt.Println(counter)
}
