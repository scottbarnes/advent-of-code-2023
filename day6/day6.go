package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type PuzzlePart int

const (
	Part1 PuzzlePart = iota
	Part2
)

// Race represents a race.
type Race struct {
	time           int
	distanceRecord int
	waysToWin      int
}

// Races represents a slice of Race.
type Races []Race

const filename = "day6_input.txt"

func main() {
	part := flag.Int("part", 1, "the integer value of the puzzle to run (1 or 2)")
	flag.Parse()
	if *part < 1 || *part > 2 {
		fmt.Println("Expected the part value to be 1 or 2.")
		os.Exit(1)
	}

	var puzzlePart PuzzlePart
	switch *part {
	case 1:
		puzzlePart = Part1
	case 2:
		puzzlePart = Part2
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	result := run(file, puzzlePart)
	fmt.Println(result)
}

// NewRace() returns a Race.
func NewRace(time int, distance int) Race {
	return Race{time: time, distanceRecord: distance}
}

// getLines() returns the file input as []string.
func getLines(file io.Reader) []string {
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// getRaces() takes a []string and returns all the Races.
func getRaces(lines []string, puzzlePart PuzzlePart) Races {
	regex := regexp.MustCompile(`\d+`)
	timeMatchStr := regex.FindAllString(lines[0], -1)
	distanceMatchStr := regex.FindAllString(lines[1], -1)

	timeMatch := []int{}
	distanceMatch := []int{}
	switch puzzlePart {
	case Part1:
		timeMatch = strSliceToIntPart1(timeMatchStr)
		distanceMatch = strSliceToIntPart1(distanceMatchStr)
	case Part2:
		timeMatch = strSliceToIntPart2(timeMatchStr)
		distanceMatch = strSliceToIntPart2(distanceMatchStr)
	}

	races := Races{}
	for i := 0; i < len(timeMatch); i++ {
		races = append(races, NewRace(timeMatch[i], distanceMatch[i]))
	}

	return races
}

func run(file io.Reader, puzzlePart PuzzlePart) int {
	lines := getLines(file)

	var races Races
	switch puzzlePart {
	case Part1:
		races = getRaces(lines, Part1)
	case Part2:
		races = getRaces(lines, Part2)
	}

	result := 1
	for _, race := range races {
		result *= race.getWaysToWin()
	}
	return result
}

// getWaysToWin() returns the number of ways to win in Part1.
func (r Race) getWaysToWin() int {
	result := 0
	for i := 0; i < r.time; i++ {
		if i*(r.time-i) > r.distanceRecord {
			result++
		}
	}

	return result
}

// Helper utilities

// strSliceToIntPart1() converts a []string to []int.
func strSliceToIntPart1(slice []string) []int {
	result := []int{}
	for _, numStr := range slice {
		num, _ := strconv.Atoi(numStr)
		result = append(result, num)
	}
	return result
}

// strSliceToIntPart2() converts a []string to []int.
func strSliceToIntPart2(slice []string) []int {
	resultStr := ""
	for _, numStr := range slice {
		resultStr += numStr
	}
	num, _ := strconv.Atoi(resultStr)
	return []int{num}
}
