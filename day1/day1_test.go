package main

import (
	"bytes"
	"testing"
)

func TestGetLineValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"two1nine", 29},
		{"eightwothree", 83},
		{"abcone2threexyz", 13},
		{"xtwone3four", 24},
		{"4nineeightseven2", 42},
		{"zoneight234", 14},
		{"7pqrstsixteen", 76},
		{"one7xctgtrtwoeightwovkv", 12},
	}

	for _, tc := range testCases {
		got := getLineValue(tc.input, []string{})
		if got != tc.expected {
			t.Errorf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestRun(t *testing.T) {
	buffer := bytes.NewBufferString("two1nine\neightwothree\nabcone2threexyz\nxtwone3four\n4nineeightseven2\nzoneight234\n7pqrstsixteen\none7xctgtrtwoeightwovkv")
	got, err := run(buffer)
	if err != nil {
		t.Error(err)
	}
	expected := 293
	if got != expected {
		t.Errorf("Expected %v, but got %v", expected, got)
	}
}

// func TestGetLineValue(t *testing.T) {
// 	testCases := []struct {
// 		input    string
// 		expected int
// 	}{
// 		{"1abc2", 12},
// 		{"pqr3stu8vwx", 38},
// 		{"a1b2c3d4e5f", 15},
// 		{"treb7uchet", 77},
// 	}

// 	for _, tc := range testCases {
// 		got := getLineValue(tc.input)
// 		if got != tc.expected {
// 			t.Errorf("Expected %v, but got %v", tc.expected, got)
// 		}
// 	}
// }

// func TestRun(t *testing.T) {
// 	buffer := bytes.NewBufferString("1abc2\npqr3stu8vwx\na1b2c3d4e5f\ntreb7uchet")
// 	got, err := run(buffer)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	expected := 142
// 	if got != expected {
// 		t.Errorf("Expected %v, but got %v", expected, got)
// 	}
// }
