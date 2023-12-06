package main

import (
	"bytes"
	"reflect"
	"testing"
)

const testInput = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

var expectedFirst = Mapping{
	name:             "seed-to-soil",
	destinationStart: 50,
	sourceStart:      98,
	range_:           2,
}

var expectedSecond = Mapping{
	name:             "seed-to-soil",
	destinationStart: 52,
	sourceStart:      50,
	range_:           48,
	// },
}

var expectedLast = Mapping{
	name:             "humidity-to-location",
	destinationStart: 56,
	sourceStart:      93,
	range_:           4,
}

func TestGetSeeds(t *testing.T) {
	expected := []int{79, 14, 55, 13}
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)
	got, err := getSeeds(lines)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %#v, but got %#v", expected, got)
	}
}

func TestNewMap(t *testing.T) {
	inputNumbers := []int{50, 98, 2}
	inputName := "seed-to-soil"
	expected := Mapping{
		name:             "seed-to-soil",
		destinationStart: 50,
		sourceStart:      98,
		range_:           2,
	}

	got := NewMapping(inputName, inputNumbers)
	if got != expected {
		t.Fatalf("Expected %#v, but got %#v", expected, got)
	}
}

func TestGetMaps(t *testing.T) {
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)

	got := getMaps(lines, ConversionMap{}, "")

	expectedLength := 18
	gotLength := len(got)
	if len(got) != expectedLength {
		t.Fatalf("Expected %d Maps, but got %d", expectedLength, gotLength)
	}

	if !reflect.DeepEqual(got[1], expectedSecond) {
		t.Fatalf("Expected %#v, but got %#v", expectedSecond, got[1])
	}

	if !reflect.DeepEqual(got[17], expectedLast) {
		t.Fatalf("Expected %#v, but got %#v", expectedLast, got[1])
	}
}

func TestDestinationInMap(t *testing.T) {
	testCases := []struct {
		name          string
		input         int
		expectedNum   int
		expectedInMap bool
	}{
		{"in map 79", 79, 81, true},
		{"not in map 14", 14, 14, false},
		{"in map 55", 55, 57, true},
		{"not in map 13", 13, 13, false},
	}

	mapping := expectedSecond

	for _, tc := range testCases {
		gotInMap, gotOutputNum := mapping.sourceInMap(tc.input)

		if gotInMap != tc.expectedInMap {
			t.Fatalf("Expected %v, but got %v for %v", tc.expectedInMap, gotInMap, tc.name)
		}

		if gotOutputNum != tc.expectedNum {
			t.Fatalf("Expected %v, but got %v for %v", tc.expectedNum, gotOutputNum, tc.name)
		}
	}
}

func TestGetDestination(t *testing.T) {
	testCases := []struct {
		source   int
		expected int
	}{
		{79, 81},
		{14, 14},
		{55, 57},
		{13, 13},
		{99, 51},
	}

	maps := ConversionMap{expectedFirst, expectedSecond, expectedLast}

	for _, tc := range testCases {
		got := maps.getDestination("seed-to-soil", tc.source)
		if got != tc.expected {
			t.Fatalf("Expected %d, but got %d", tc.expected, got)
		}
	}
}

func TestGetSeedLocation(t *testing.T) {
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)
	maps := getMaps(lines, ConversionMap{}, "")
	// allMaps := Maps{maps}

	testCases := []struct {
		seed     int
		expected int
	}{
		{79, 82},
		{14, 43},
		{55, 86},
		{13, 35},
	}

	for _, tc := range testCases {
		got := getSeedLocation(tc.seed, conversionMapNames, maps)
		if got != tc.expected {
			t.Fatalf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestGetSeedRanges(t *testing.T) {
	rawRangeStr := []string{"48", "53", "200", "2", "10", "2"}
	expected := []SourceRange{
		{48, 100, 53},
		{200, 201, 2},
		{10, 11, 2},
	}

	got := getSeedRanges(rawRangeStr, []SourceRange{})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestGetSeedDestinations(t *testing.T) {
	map1 := NewMapping("seed-to-soil", []int{50, 110, 2}) // everything passes through this.
	map2 := NewMapping("seed-to-soil", []int{52, 50, 48}) // 50-97
	seedRange1 := SourceRange{48, 100, 53}                // 2 below, middle 3 above
	seedRange2 := SourceRange{200, 201, 2}                // above everything
	seedRange3 := SourceRange{10, 11, 2}                  // below everything
	seedRange4 := SourceRange{45, 54, 10}                 // Part below, part in.
	seedRange5 := SourceRange{95, 104, 10}
	seedRange6 := SourceRange{49, 49, 1} // off by one
	seedRange7 := SourceRange{50, 50, 1} // off by one
	seedRange8 := SourceRange{51, 51, 1} // off by one
	seedRange9 := SourceRange{96, 96, 1}
	seedRange10 := SourceRange{97, 97, 1}
	seedRange11 := SourceRange{98, 98, 1}
	expected := []SourceRange{
		{48, 49, 0},  // intersection starts at 50
		{98, 100, 0}, // intersection ends at 97
		{200, 201, 2},
		{10, 11, 2},
		{45, 49, 0},
		{98, 104, 0},
		{49, 49, 1},
		{98, 98, 1},
		{52, 99, 0}, // mapped because it's an intersection with a map
		{52, 56, 0},
		{97, 99, 0},
		{52, 52, 0},
		{53, 53, 0},
		{98, 98, 0},
		{99, 99, 0},
	}

	seedRanges := []SourceRange{seedRange1, seedRange2, seedRange3, seedRange4, seedRange5, seedRange6, seedRange7, seedRange8, seedRange9, seedRange10, seedRange11}
	conversionMapsSmall := []string{"seed-to-soil"}
	maps := ConversionMap{map1, map2}

	got := getFinalLocations(seedRanges, conversionMapsSmall, maps)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestLowestSeedInRange(t *testing.T) {
	expected := 46
	buffer := bytes.NewBufferString(testInput)
	got, err := run(buffer, Part2)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}
