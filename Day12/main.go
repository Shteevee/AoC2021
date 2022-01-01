package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

type Cave struct {
	isSmall bool
	paths   []*Cave
	name    string
}

func caveExists(uniqueCaves []*Cave, node string) bool {
	for _, uniqueNode := range uniqueCaves {
		if uniqueNode.name == node {
			return true
		}
	}

	return false
}

func isLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func contains(path []string, caveName string) bool {
	for _, cave := range path {
		if cave == caveName {
			return true
		}
	}

	return false
}

func findLinkedCaves(cave *Cave, paths [][]string) []string {
	linkedPaths := make([]string, 0)
	for _, path := range paths {
		if contains(path, cave.name) {
			for _, caveName := range path {
				if caveName != cave.name {
					linkedPaths = append(linkedPaths, caveName)
				}
			}
		}
	}

	return linkedPaths
}

func findCave(caves []*Cave, linkedCaveName string) (*Cave, bool) {
	for _, cave := range caves {
		if cave.name == linkedCaveName {
			return cave, true
		}
	}

	return &Cave{}, false
}

func createCavePaths(cave *Cave, uniqueCaves []*Cave, paths [][]string) []*Cave {
	cavePaths := make([]*Cave, 0)
	linkedPaths := findLinkedCaves(cave, paths)

	if cave.name == "end" {
		return cavePaths
	}

	for _, path := range linkedPaths {
		linkedCave, caveExists := findCave(uniqueCaves, path)
		if caveExists {
			cavePaths = append(cavePaths, linkedCave)
		}
	}

	return cavePaths
}

func parseCaveSystem(scanner *bufio.Scanner) Cave {
	paths := make([][]string, 0)
	for scanner.Scan() {
		paths = append(paths, strings.Split(scanner.Text(), "-"))
	}

	uniqueCaves := make([]*Cave, 0)
	for _, path := range paths {
		for _, node := range path {
			if !caveExists(uniqueCaves, node) {
				uniqueCaves = append(
					uniqueCaves,
					&Cave{
						isSmall: isLower(node),
						name:    node,
						paths:   make([]*Cave, 0),
					},
				)
			}
		}
	}

	for _, cave := range uniqueCaves {
		cave.paths = createCavePaths(cave, uniqueCaves, paths)
	}

	startCave, _ := findCave(uniqueCaves, "start")

	return *startCave
}

//there is almost certainly a better way to check this
//it hammers the performance
func smallCaveVisitedTwice(currentPath []string) string {
	for i := range currentPath {
		for j := range currentPath {
			if isLower(currentPath[i]) && currentPath[i] == currentPath[j] && i != j {
				return currentPath[i]
			}
		}
	}

	return ""
}

func explored(pathOption Cave, currentPath []string) bool {
	visitedHere := pathOption.isSmall && contains(currentPath, pathOption.name) && smallCaveVisitedTwice(currentPath) != ""
	for _, prevCave := range currentPath {
		if prevCave == pathOption.name {
			if prevCave == "start" || prevCave == "end" || visitedHere {
				return true
			}
		}
	}

	return false
}

func explore(cave Cave, currentPath []string) [][]string {
	currentPath = append(currentPath, cave.name)

	if cave.name == "end" {
		initPathMap := make([][]string, 1)
		initPathMap[0] = currentPath
		return initPathMap
	}

	exploredPaths := make([][]string, 0)
	for _, pathOption := range cave.paths {
		if !explored(*pathOption, currentPath) {
			exploredPaths = append(exploredPaths, explore(*pathOption, currentPath)...)
		}
	}

	if len(exploredPaths) > 0 {
		return exploredPaths
	}

	return make([][]string, 0)
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
	startCave := parseCaveSystem(scanner)

	routes := explore(startCave, make([]string, 0))

	elapsed := time.Since(start)
	fmt.Println(len(routes))
	log.Printf("Time taken: %s", elapsed)
}
