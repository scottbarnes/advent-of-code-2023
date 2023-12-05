package main

import (
	"bytes"
	"reflect"
	"testing"
)

const (
	testCards      = "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53\nCard 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19\nCard 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1\nCard 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83\nCard 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36\nCard 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11"
	singleTestCard = "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53"
)

func TestNewCard(t *testing.T) {
	deck := NewCardDeck()
	expected := Card{
		deck:           deck,
		cardNo:         1,
		copies:         1,
		matches:        4,
		winningNumbers: []int{41, 48, 83, 86, 17},
		yourNumbers:    []int{83, 86, 6, 31, 17, 9, 48, 53},
	}
	got := NewCard(deck, singleTestCard, 0)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestCardGetPoints(t *testing.T) {
	deck := NewCardDeck()
	card := NewCard(deck, singleTestCard, 0)
	expected := 8
	got := card.getPoints()
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}

func TestRun(t *testing.T) {
	buffer := bytes.NewBufferString(testCards)
	got := run(buffer, Points)
	expected := 13
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}

func TestRunCopies(t *testing.T) {
	buffer := bytes.NewBufferString(testCards)
	got := run(buffer, Copies)
	expected := 30
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}
