package main

import (
	"bytes"
	"reflect"
	"testing"
)

var testInput = `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

var (
	hand1   = newHand([]int{3, 2, 10, 3, 13}, 765)    // One pair       = 4
	hand2   = newHand([]int{10, 5, 5, 11, 5}, 684)    // Three of a kind
	hand3   = newHand([]int{13, 13, 6, 7, 7}, 28)     // Two pair
	hand4   = newHand([]int{13, 10, 11, 11, 10}, 220) // Two pair       = 3
	hand5   = newHand([]int{12, 12, 12, 11, 14}, 483) // Three of a kind = 3
	hand6   = newHand([]int{1, 1, 1, 2, 2}, 200)      // Full house     = 2
	hand7   = newHand([]int{2, 2, 2, 2, 3}, 400)      // Four of a kind = 2
	hand8   = newHand([]int{3, 3, 3, 3, 3}, 500)      // Five of a kind = 1
	hand9   = newHand([]int{1, 2, 3, 4, 5}, 1)        // High card      = 5
	hand1p2 = newHand([]int{3, 2, 10, 3, 12}, 765)    // One pair       = 4
	hand2p2 = newHand([]int{10, 5, 5, 1, 5}, 684)     // Four of a kind
	hand3p2 = newHand([]int{12, 12, 6, 7, 7}, 28)     // Two pair
	hand4p2 = newHand([]int{12, 10, 1, 1, 10}, 220)   // Four of a kind
	hand5p2 = newHand([]int{11, 11, 11, 1, 13}, 483)  // Four of a kind
	hand6p2 = newHand([]int{1, 1, 1, 2, 2}, 400)      // Five of a kind
	hand7p2 = newHand([]int{2, 2, 2, 2, 3}, 400)      // Five of a kind = 2
	hand8p2 = newHand([]int{1, 1, 1, 1, 2}, 500)      // Five of a kind = 1
)

func TestGetHandsPart1(t *testing.T) {
	expected := Hands{hand1, hand2, hand3, hand4, hand5}
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)
	got := getHands(lines, Part1)

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestGetHandsPart2(t *testing.T) {
	expected := Hands{hand1p2, hand2p2, hand3p2, hand4p2, hand5p2}
	buffer := bytes.NewBufferString(testInput)
	lines := getLines(buffer)
	got := getHands(lines, Part2)

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestClassifyHandPart1(t *testing.T) {
	testCases := []struct {
		hand     Hand
		expected int
	}{
		{hand1, OnePair},
		{hand2, ThreeOfAKind},
		{hand3, TwoPair},
		{hand4, TwoPair},
		{hand5, ThreeOfAKind},
		{hand6, FullHouse},
		{hand7, FourOfAKind},
		{hand8, FiveOfAKind},
		{hand9, HighCard},
	}

	for _, tc := range testCases {
		got := tc.hand.classifyHand()
		if got != tc.expected {
			t.Fatalf("Expected %v, but got %v", tc.expected, got)
		}
	}
}

func TestWildCards(t *testing.T) {
	testCases := []struct {
		hand     Hand
		expected int
	}{
		{hand1p2, OnePair},
		{hand2p2, FourOfAKind},
		{hand3p2, TwoPair},
		{hand4p2, FourOfAKind},
		{hand5p2, FourOfAKind},
		{hand6p2, FiveOfAKind},
		{hand7p2, FourOfAKind},
		{hand8p2, FiveOfAKind},
	}

	for _, tc := range testCases {
		tc.hand.makeWildCards()
		got := tc.hand.classifyHand()
		if got != tc.expected {
			t.Fatalf("Expected %v, but got %v for Hand: %v", tc.expected, got, tc.hand)
		}
	}
}

func TestIsStrongerThan(t *testing.T) {
	testCases := []struct {
		name       string
		firstHand  Hand
		secondHand Hand
		expected   bool
	}{
		{"hand1 is not stronger than hand2", hand1, hand2, false},
		{"hand2 is stronger than hand1", hand2, hand1, true},
		{"hand3 is stronger than hand4", hand3, hand4, true},
		{"hand4 is not stronger than hand 3", hand4, hand3, false},
	}

	for _, tc := range testCases {
		got := tc.firstHand.isStrongerThan(tc.secondHand)
		if got != tc.expected {
			t.Fatalf("Expected %v, but got %v: %v", tc.expected, got, tc.name)
		}
	}
}

func TestSortHandsByRank(t *testing.T) {
	expected := Hands{hand1, hand4, hand3, hand2, hand5}
	hands := Hands{hand1, hand2, hand3, hand4, hand5}
	hands.sort()
	if !reflect.DeepEqual(hands, expected) {
		t.Fatalf("Expected %v, but got %v", expected, hands)
	}
}

func TestGetTotalWinnings(t *testing.T) {
	expected := 6440
	// sortedHands := Hands{hand1, hand4, hand3, hand2, hand5}
	hands := Hands{hand1, hand2, hand3, hand4, hand5}
	got := hands.getTotalWinnings()
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}

func TestRunPart1(t *testing.T) {
	buffer := bytes.NewBufferString(testInput)
	expected := 6440
	got := run(buffer, Part1)
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}

func TestRunPart2(t *testing.T) {
	buffer := bytes.NewBufferString(testInput)
	expected := 5905
	got := run(buffer, Part2)
	if got != expected {
		t.Fatalf("Expected %d, but got %d", expected, got)
	}
}
