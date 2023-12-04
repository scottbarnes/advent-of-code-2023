package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const filename = "day3_input.txt"

var (
	indexRegexNumber   = regexp.MustCompile(`\d+`)
	indexRegexAsterisk = regexp.MustCompile(`\*`)
	indexRegexSymbol   = regexp.MustCompile(`[^\d|^.|^\s]`)
)

type FindType int

const (
	Gears FindType = iota
	PartNumbers
)

type MatchIndicies struct {
	number int
	start  int
	end    int
}

// PartNumber represents a number in the schematic.
// `number` is the literal number and `indicies` are the indicies of each
// digit, including +1 on either side to aid in finding adjacent symbols.
type PartNumber struct {
	number  int
	indices []int
}

// NewPartNumber creates a new PartNumber.
func NewPartNumber(numStr string, start int, end int) PartNumber {
	number, _ := strconv.Atoi(numStr)
	indicies := makeRange(start-1, end+1) // Expand range to include adjacent indices.
	return PartNumber{
		number:  number,
		indices: indicies,
	}
}

// isAdjacent checks if a PartNumber is adjacent to a given index.
// The adjacent check includes numbers that are diagonal from an index.
func (pn PartNumber) isAdjacent(index int) bool {
	for _, number := range pn.indices {
		if number == index {
			return true
		}
	}

	return false
}

// main is the program entrypoint and accepts two args: 'partnumbers' or 'gears'.
// This program identifies part numbers (i.e. numbers adjacent to symbols)
// returns their sum with the 'partnumbers' argument, or the sum of gear ratios
// (i.e. the sum of the result of multiplying two and only two numbers adjacent
// to asterisk) with the 'gears' argument.
func main() {
	args := os.Args
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer file.Close()

	switch args[1] {
	case "partnumbers":
		result, err := readSchematic(file, PartNumbers)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		fmt.Println(result)
	case "gears":
		result, err := readSchematic(file, Gears)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		fmt.Println(result)
	default:
		fmt.Println("Expected argument to be 'partnumbers' or 'gears'")
	}
}

// Generate a range between two numbers, excluding the last number.
func makeRange(start, end int) []int {
	var result []int
	for i := start; i < end; i++ {
		result = append(result, i)
	}
	return result
}

// Get the relevant symbol matches for a line.
func getSymbolMatches(findType FindType, line string) [][]int {
	switch findType {
	case PartNumbers:
		return indexRegexSymbol.FindAllStringIndex(line, -1)
	case Gears:
		return indexRegexAsterisk.FindAllStringIndex(line, -1)
	default:
		fmt.Println("Invalid findType: must be PartNumbers or Gears.")
		return nil
	}
}

// calculateSum adds up all the values of a []PartNumber per the rules of the FindType.
// PartNumbers adjacent to symbols have their numerical value added to the total sum.
// Gears have their two part numbers multiplied then added to the total sum.
func calculateSum(findType FindType, candidates []PartNumber, index int) int {
	var adjacentNumbers []int
	for _, candidate := range candidates {
		if candidate.isAdjacent(index) {
			adjacentNumbers = append(adjacentNumbers, candidate.number)
		}
	}

	var sum int
	switch findType {
	case PartNumbers:
		for _, number := range adjacentNumbers {
			sum += number
		}
	case Gears:
		if len(adjacentNumbers) == 2 {
			sum += adjacentNumbers[0] * adjacentNumbers[1]
		}
	default:
		fmt.Println("Invalid findType: must be PartNumbers or Gears.")
		return 0
	}

	return sum
}

// readSchematic reads through a schematic and adds up numbers per the rules
// for part numbers and gears.
func readSchematic(reader io.Reader, findType FindType) (int, error) {
	scanner := bufio.NewScanner(reader)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Process all symbols on all lines and get the sum per the findType.
	total := 0
	for lineIndex, line := range lines {
		lineMatchIndices := getSymbolMatches(PartNumbers, line)
		for _, match := range lineMatchIndices {
			candidates := getAllCandidates(lineIndex, lines)
			total += calculateSum(findType, candidates, match[0])
		}
	}

	return total, nil
}

// makeLineCandidates creates a []PartNumber of every PartNumber on a line.
func makeLineCandidates(line string) []PartNumber {
	matchNumbers := indexRegexNumber.FindAllString(line, -1)      // The number itself.
	matchIndices := indexRegexNumber.FindAllStringIndex(line, -1) // The number's start/end index.
	var candidates []PartNumber
	for idx, numStr := range matchNumbers {
		partNumber := NewPartNumber(numStr, matchIndices[idx][0], matchIndices[idx][1])
		candidates = append(candidates, partNumber)
	}
	return candidates
}

// getAllCandidates creates a []PartNumber of every number candidate
// Look for a matching index on the line above, below, and the same line as
// the symbol.
func getAllCandidates(baseLineNumber int, lines []string) []PartNumber {
	var candidates []PartNumber
	// Above line
	if baseLineNumber > 0 {
		candidates = append(candidates, makeLineCandidates(lines[baseLineNumber-1])...)
	}

	// Below line
	if baseLineNumber < len(lines)-1 {
		candidates = append(candidates, makeLineCandidates(lines[baseLineNumber+1])...)
	}

	// Same line
	candidates = append(candidates, makeLineCandidates(lines[baseLineNumber])...)
	return candidates
}
