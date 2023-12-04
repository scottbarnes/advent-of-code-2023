package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var numberMap = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
	"1":     "1",
	"2":     "2",
	"3":     "3",
	"4":     "4",
	"5":     "5",
	"6":     "6",
	"7":     "7",
	"8":     "8",
	"9":     "9",
	"0":     "0",
}

var inputFile = "day1_input.txt"

// main prints the total calibration values in an input file.
func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	result, err := run(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

// getLineValue returns the first and last calibration values from a line by
// recursively advancing one character and looking for a matching numbers, then
// taking the first and last match.
// E.g., "one7xctgtrtwoeightwovkv" would return 12.
func getLineValue(line string, acc []string) int {
	if len(line) == 0 {
		result, _ := strconv.Atoi(acc[0] + acc[len(acc)-1])
		return result
	}

	for key, value := range numberMap {
		if strings.HasPrefix(line, key) {
			acc = append(acc, value)
		}
	}

	return getLineValue(line[1:], acc)
}

// run reads through calibration lines and returns their sum or an error.
func run(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	total := 0
	for scanner.Scan() {
		total += getLineValue(scanner.Text(), []string{})
	}

	return total, nil
}

// func main() {
// 	file, err := os.Open("day1_input.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := run(file)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(result)
// }

// func getLineValue(line string) int {
// 	numbers := []rune{}
// 	for _, char := range line {
// 		if unicode.IsDigit(char) {
// 			numbers = append(numbers, char)
// 		}
// 	}

// 	firstLast := []rune{numbers[0], numbers[len(numbers)-1]}
// 	result, _ := strconv.Atoi(string(firstLast))

// 	return result
// }

// func run(reader io.Reader) (int, error) {
// 	scanner := bufio.NewScanner(reader)
// 	total := 0
// 	for scanner.Scan() {
// 		total += getLineValue(scanner.Text())
// 	}

// 	return total, nil
// }
