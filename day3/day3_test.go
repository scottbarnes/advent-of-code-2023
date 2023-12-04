package main

import (
	"bytes"
	"testing"
)

const testSchematic = "467..114..\n...*......\n..35..633.\n......#...\n617*......\n.....+..58\n..592.....\n......755.\n...$.*....\n.664.598.."

func TestReadSchematicPartOne(t *testing.T) {
	buffer := bytes.NewBufferString(testSchematic)
	got, err := readSchematic(buffer, PartNumbers)
	if err != nil {
		t.Error(err)
	}

	expected := 4361
	if got != expected {
		t.Errorf("Expected %d, but got %d", expected, got)
	}
}

func TestReadSchematicPartTwo(t *testing.T) {
	buffer := bytes.NewBufferString(testSchematic)
	got, err := readSchematic(buffer, Gears)
	if err != nil {
		t.Error(err)
	}

	expected := 467835
	if got != expected {
		t.Errorf("Expected %d, but got %d", expected, got)
	}
}
