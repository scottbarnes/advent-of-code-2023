package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
)

const filename = "day5_input.txt"

var (
	regexSeeds         = regexp.MustCompile(`(\d+)`)
	conversionMapNames = []string{
		"seed-to-soil",
		"soil-to-fertilizer",
		"fertilizer-to-water",
		"water-to-light",
		"light-to-temperature",
		"temperature-to-humidity",
		"humidity-to-location",
	}
)

type FindType int

const (
	Part1 FindType = iota
	Part2
)

// Mapping is an individual map of the source to destination.
type Mapping struct {
	name             string
	destinationStart int
	sourceStart      int
	range_           int
}

// ConversionMap is a collection of transformation mappings, such as seed-to-soil.
type ConversionMap []Mapping

// SourceRange represents a range of numbers for use with the source of a map.
type SourceRange struct {
	start  int
	end    int
	range_ int
}

func main() {
	part := flag.Int("part", 2, "The puzzle part to run (i.e. either 1 or 2)")
	flag.Parse()
	if *part < 1 || *part > 2 {
		fmt.Printf("Expected argument to -part: '1' or '2'.")
		os.Exit(1)
	}

	var puzzlePart FindType
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

	result, err := run(file, puzzlePart)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// run() is the entrypoint to the program and contains the main logic.
func run(file io.Reader, puzzlePart FindType) (int, error) {
	lines := getLines(file)
	maps := getMaps(lines, ConversionMap{}, "")
	locations := []int{}
	lowest := math.MaxInt64

	switch puzzlePart {
	case Part1:
		seeds, err := getSeeds(lines)
		if err != nil {
			return 0, err
		}
		for _, seed := range seeds {
			locations = append(locations, getSeedLocation(seed, conversionMapNames, maps))
		}
		for _, num := range locations {
			if num < lowest {
				lowest = num
			}
		}
	case Part2:
		matches := regexSeeds.FindAllString(lines[0], -1)
		if matches == nil {
			return 0, errors.New("can't find seed")
		}
		seedRanges := getSeedRanges(matches, []SourceRange{})
		locations := getFinalLocations(seedRanges, conversionMapNames, maps)
		for _, location := range locations {
			if location.end < lowest {
				lowest = location.start
			}
		}
	}

	return lowest, nil
}

// Shared code between Part 1 and Part 2.
// getLines() gets the lines from a file as []string.
func getLines(file io.Reader) []string {
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// getMaps returns a Mappings representing all the destination-source maps.
func getMaps(lines []string, acc ConversionMap, lastMapName string) ConversionMap {
	if len(lines) == 0 {
		return acc
	}

	regexpMapName := regexp.MustCompile(`(\w+.*) (?:map:)`)
	regexpMap := regexp.MustCompile(`^(\d+)\s+(\d+)\s+(\d+)$`)
	line := lines[0]
	matchMapName := regexpMapName.FindStringSubmatch(line)
	matchMap := regexpMap.FindStringSubmatch(line)

	if matchMapName != nil {
		lastMapName = matchMapName[1:][0]
	} else if matchMap != nil {
		slice := []int{}
		for _, numStr := range matchMap[1:] {
			num, _ := strconv.Atoi(numStr)
			slice = append(slice, num)
		}
		mapping := NewMapping(lastMapName, slice)
		acc = append(acc, mapping)
	}
	return getMaps(lines[1:], acc, lastMapName)
}

// Part 2 code.

// convertRanges() is a recursive function to process a group of map directives and
// returns the a []sourceRange of the final locations, whether a source was
// mapped to a new destination, or whether it fell through in Part 2.
func (tm ConversionMap) convertRanges(unProcessed []SourceRange, converted []SourceRange) []SourceRange {
	if len(tm) == 0 {
		return append(unProcessed, converted...)
	}

	mapPart := tm[0]
	fellThrough := []SourceRange{}

	for _, seedRange := range unProcessed {
		mapSourceEnd := mapPart.sourceStart + mapPart.range_ - 1

		// Pass through ranges lacking an intersection with this map part's source.
		if seedRange.end < mapPart.sourceStart || seedRange.start > mapSourceEnd {
			fellThrough = append(fellThrough, seedRange)
			continue
		}

		// Find the start and end of the intersection so we know the range that
		// must shift via the map source.
		intersectionStart := max(seedRange.start, mapPart.sourceStart)
		intersectionEnd := min(seedRange.end, mapSourceEnd)

		// If numbers are above or below the source's intersection with the range,
		// then pass them through to the next map.
		if intersectionStart > seedRange.start {
			below := SourceRange{start: seedRange.start, end: intersectionStart - 1}
			fellThrough = append(fellThrough, below)
		}

		if intersectionEnd < seedRange.end {
			above := SourceRange{start: intersectionEnd + 1, end: seedRange.end}
			fellThrough = append(fellThrough, above)
		}

		// Calculate the shifted intersection range based on the map part's destination.
		offset := intersectionStart - mapPart.sourceStart
		intersectionLength := intersectionEnd - intersectionStart
		shiftedIntersection := SourceRange{start: mapPart.destinationStart + offset, end: mapPart.destinationStart + offset + intersectionLength}
		converted = append(converted, shiftedIntersection)
	}

	return tm[1:].convertRanges(fellThrough, converted)
}

// getFinalLocations() is a recursive function that processes map collections
// one at a time and returns a []sourceRange of the final location.
// E.g. first it processes the locations for seed-to-soil, then soil-to-fertilizer, etc.
func getFinalLocations(seeds []SourceRange, conversionMaps []string, maps ConversionMap) []SourceRange {
	if len(conversionMaps) == 0 {
		return seeds
	}

	// Gather the maps components (e.g. all `seed-to-soil` maps) for this round.
	conversionMap := ConversionMap{}
	for _, map_ := range maps {
		if map_.name == conversionMaps[0] {
			conversionMap = append(conversionMap, map_)
		}
	}

	seedRanges := conversionMap.convertRanges(seeds, []SourceRange{})
	return getFinalLocations(seedRanges, conversionMaps[1:], maps)
}

// getSeedRanges() is a takes a []string of start/end range pairs and converts
// them into []SourceRange for Part 2.
func getSeedRanges(rawRangeStr []string, seeds []SourceRange) []SourceRange {
	if len(rawRangeStr) == 0 {
		return seeds
	}

	startStr, rangeStr := rawRangeStr[0], rawRangeStr[1]
	start, _ := strconv.Atoi(startStr)
	range_, _ := strconv.Atoi(rangeStr)
	end := start + range_ - 1

	seeds = append(seeds, SourceRange{start: start, end: end, range_: range_})
	return getSeedRanges(rawRangeStr[2:], seeds)
}

// NewMapping() creates a new Mapping object.
func NewMapping(name string, numbers []int) Mapping {
	destinationStart := numbers[0]
	sourceStart := numbers[1]
	range_ := numbers[2]
	return Mapping{name, destinationStart, sourceStart, range_}
}

// Part 1 code

// getSeeds() gets seeds for Part 1.
func getSeeds(lines []string) ([]int, error) {
	matches := regexSeeds.FindAllString(lines[0], -1)
	if matches == nil {
		return []int{}, errors.New("can't find seed")
	}
	seeds := []int{}
	for _, numStr := range matches {
		seed, _ := strconv.Atoi(numStr)
		seeds = append(seeds, seed)
	}

	return seeds, nil
}

// sourceInMap() returns true if a source is in the map, and also returns
// its destination for Part 1.
func (m *Mapping) sourceInMap(source int) (bool, int) {
	if m.sourceStart <= source && source <= m.sourceStart+m.range_ {
		return true, m.destinationStart + (source - m.sourceStart)
	}

	return false, source
}

// getDestination() gets a seed destination from a map source in Part 1.
func (cMap ConversionMap) getDestination(mapNamp string, source int) int {
	sourceCopy := source
	for _, mapping := range cMap {
		if mapping.name == mapNamp {
			inMap, num := mapping.sourceInMap(sourceCopy)
			if inMap {
				return num
			}
		}
	}

	return sourceCopy
}

// getSeedLocation() gets the final seed location in Part 1.
func getSeedLocation(seed int, conversionMapNames []string, conversionMap ConversionMap) int {
	if len(conversionMapNames) == 0 {
		return seed
	}

	mapName := conversionMapNames[0]
	result := conversionMap.getDestination(mapName, seed)
	return getSeedLocation(result, conversionMapNames[1:], conversionMap)
}

// Utility functions
func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
