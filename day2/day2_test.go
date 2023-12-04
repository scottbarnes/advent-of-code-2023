package main

import (
	"bytes"
	"testing"
)

func TestGetGameTotals(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			"Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
			1,
		},
		{
			"Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
			2,
		},
		{
			"Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
			0,
		},
		{
			"Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
			0,
		},
		{
			"Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
			5,
		},
		{
			"Game 33: 4 red; 3 red; 2 red, 1 green, 1 blue; 1 green; 1 blue, 1 red",
			33,
		},
		{
			"Game 70: 12 green, 1 blue, 4 red; 8 green, 1 red; 1 blue, 8 green; 2 green, 3 red; 5 green, 4 red; 2 blue, 12 green, 1 red",
			70,
		},
	}

	for _, tc := range testCases {
		game := parseGame(tc.input)
		got, err := calculateGameValue(game)
		if err != nil {
			t.Error(err)
		}

		if got != tc.expected {
			t.Errorf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestGetCubeTotals(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			"Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
			48,
		},
		{
			"Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
			12,
		},
		{
			"Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
			1560,
		},
		{
			"Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
			630,
		},
		{
			"Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
			36,
		},
	}

	for _, tc := range testCases {
		game := parseGame(tc.input)
		got, err := calculateCubePower(game)
		if err != nil {
			t.Error(err)
		}

		if got != tc.expected {
			t.Errorf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestRunGameTotals(t *testing.T) {
	buffer := bytes.NewBufferString("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green\nGame 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue\nGame 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red\nGame 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red\nGame 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green\nGame 33: 4 red; 3 red; 2 red, 1 green, 1 blue; 1 green; 1 blue, 1 red\nGame 70: 12 green, 1 blue, 4 red; 8 green, 1 red; 1 blue, 8 green; 2 green, 3 red; 5 green, 4 red; 2 blue, 12 green, 1 red")
	got, err := processGames(buffer, GameTotals)
	if err != nil {
		t.Error(err)
	}

	expected := 111
	if got != expected {
		t.Errorf("Expected %v, but got %v", expected, got)
	}
}

func TestRunCubeTotals(t *testing.T) {
	buffer := bytes.NewBufferString("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green\nGame 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue\nGame 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red\nGame 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red\nGame 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green")
	got, err := processGames(buffer, CubeTotals)
	if err != nil {
		t.Error(err)
	}

	expected := 2286
	if got != expected {
		t.Errorf("Expected %v, but got %v", expected, got)
	}
}
