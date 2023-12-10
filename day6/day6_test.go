package main

import (
	"bytes"
	"reflect"
	"testing"
)

var testInput = `Time:      7  15   30
Distance:  9  40  200`

var testInput2 = `Time:      71530
Distance:  940200`

var (
	race1 = NewRace(7, 9)
	race2 = NewRace(15, 40)
	race3 = NewRace(30, 200)
)

var races = Races{race1, race2, race3}

func TestGetRaces(t *testing.T) {
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)
	expected := races

	got := getRaces(lines, Part1)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestGetWaysToWin(t *testing.T) {
	testCases := []struct {
		race     Race
		expected int
	}{
		{race1, 4},
		{race2, 8},
		{race3, 9},
	}

	for _, tc := range testCases {
		got := tc.race.getWaysToWin()
		if got != tc.expected {
			t.Fatalf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestRunPart1(t *testing.T) {
	expected := 288
	buffer := bytes.NewBufferString(testInput)
	got := run(buffer, Part1)
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}

func TestRunPart2(t *testing.T) {
	expected := 71503
	buffer := bytes.NewBufferString(testInput)
	got := run(buffer, Part2)
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}
