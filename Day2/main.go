package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

	distance := 0
	depth := 0
	aim := 0

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " ")

		direction := split[0]
		magnitude, _ := strconv.Atoi(split[1])

		switch direction {
		case "forward":
			distance = distance + magnitude
			depth = depth + (aim * magnitude)
		case "up":
			aim = aim - magnitude
		case "down":
			aim = aim + magnitude
		}
	}

	fmt.Println("Distance: ", distance)
	fmt.Println("Depth: ", depth)
	fmt.Println("Result: ", (distance * depth))
}
