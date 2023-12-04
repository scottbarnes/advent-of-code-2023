package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const filename = "day2_input.txt"

var (
	bagMap     = map[string]int{"red": 12, "blue": 14, "green": 13}
	cubesRegex = regexp.MustCompile(`(?:\d+ \w+)`)
	gameRegex  = regexp.MustCompile(`Game (\d+)`)
)

type Cubes struct {
	color  string
	number int
}

type Game struct {
	number int
	cubes  []Cubes
}

type GameType int

const (
	GameTotals GameType = iota
	CubeTotals
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Invalid argument. Expected 'gametotals' or 'cubetotals'.")
		os.Exit(0)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var result int
	switch args[0] {
	case "gametotals":
		result, err = processGames(file, GameTotals)
	case "cubetotals":
		result, err = processGames(file, CubeTotals)
	default:
		fmt.Println("Invalid argument. Expected 'gametotals' or 'cubetotals'.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error processing games: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func processGames(reader io.Reader, gameType GameType) (int, error) {
	scanner := bufio.NewScanner(reader)
	result := 0
	for scanner.Scan() {
		game := parseGame(scanner.Text())
		var gameValue int
		var err error
		switch gameType {
		case GameTotals:
			gameValue, err = calculateGameValue(game)
		case CubeTotals:
			gameValue, err = calculateCubePower(game)
		default:
			return 0, fmt.Errorf("unknown game type")
		}
		if err != nil {
			return 0, err
		}
		result += gameValue
	}

	return result, nil
}

// parseGame parses a game line and returns a Game struct.
func parseGame(line string) Game {
	cubesMatch := cubesRegex.FindAllStringSubmatch(line, -1)

	var cubes []Cubes
	for _, match := range cubesMatch {
		numAndColor := strings.Split(match[0], " ")
		color := numAndColor[1]
		num, _ := strconv.Atoi(numAndColor[0])
		cubes = append(cubes, Cubes{color, num})
	}

	gameNum, _ := strconv.Atoi(gameRegex.FindStringSubmatch(line)[1])
	return Game{number: gameNum, cubes: cubes}
}

func calculateGameValue(game Game) (int, error) {
	for _, cubes := range game.cubes {
		if cubes.number > bagMap[cubes.color] {
			return 0, nil
		}
	}

	return game.number, nil
}

// calculateCubePower takes, for each cube color, the highest number of cubes,
// multiplies them, and returns the result.
func calculateCubePower(game Game) (int, error) {
	// Find the highest values for each color.
	colorsMax := map[string]int{"red": math.MinInt64, "blue": math.MinInt64, "green": math.MinInt64}
	for _, cubes := range game.cubes {
		if cubes.number > colorsMax[cubes.color] {
			colorsMax[cubes.color] = cubes.number
		}
	}

	// Multiply max color values
	result := 1
	for _, value := range colorsMax {
		result *= value
	}
	return result, nil
}
