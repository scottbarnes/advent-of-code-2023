package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const filename = "day7_input.txt"

type Hand struct {
	cards     []int
	bid       int
	wildCards []int
}

type Hands []Hand

// Part 1 cards
const (
	T1 = iota + 10
	J1
	Q1
	K1
	A1
)

// Part 2 cards.
const J2 = 1

const (
	T2 = iota + 10
	Q2
	K2
	A2
)

// Represent the hand types in rank order, with the first being the lowest.
const (
	HighCard = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type PuzzlePart int

const (
	Part1 PuzzlePart = iota
	Part2
)

// newHand() creates a new Hand of cards.
func newHand(cards []int, bid int) Hand {
	return Hand{cards, bid, []int{}}
}

// h.makeWildHands() will set .wildCards for each hand, if applicable.
func (h *Hands) makeWildHands() {
	for i := range *h {
		(*h)[i].makeWildCards()
	}
}

// h.makeWildCards() uses any wild cards ("J") to create a hand of the highest
// possible value.
func (h *Hand) makeWildCards() {
	// Get the indexes of Js to use for generation and replacement.
	jIndexes := []int{}
	for i, v := range h.cards {
		if v == J2 {
			jIndexes = append(jIndexes, i)
		}
	}

	// Exit if no Jokers.
	if len(jIndexes) <= 0 {
		return
	}

	// Generate all the possible combinations for the number of Js.
	jokerCombos := getCombinations(jIndexes)

	// Slot the generated combinations into the proper index, and check if
	// it has generated a new highestHand.
	highestHand := newHand(h.cards, 0)
	for _, combo := range jokerCombos {
		tmpCards := make([]int, len(highestHand.cards))
		copy(tmpCards, highestHand.cards)
		for jIdx, jokerIdxInHand := range jIndexes {
			tmpCards[jokerIdxInHand] = combo[jIdx]
		}
		tmpHand := newHand(tmpCards, 0)
		if tmpHand.classifyHand() > highestHand.classifyHand() {
			highestHand = tmpHand
		}
	}

	h.wildCards = highestHand.cards
}

func main() {
	part := flag.Int("part", 1, "specify a puzzle part (1 or 2)")
	flag.Parse()
	if *part < 1 || *part > 2 {
		fmt.Println("Expected 'part' argument to be 1 or 2")
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	result := 0
	switch *part {
	case 1:
		result = run(file, Part1)
	case 2:
		result = run(file, Part2)
	}

	fmt.Println("result:", result)
}

// run() is the entrypoint to the program.
func run(file io.Reader, part PuzzlePart) int {
	lines := getLines(file)
	hands := getHands(lines, part)

	switch part {
	case Part1:
		return hands.getTotalWinnings()
	case Part2:
		hands.makeWildHands()
		return hands.getTotalWinnings()
	}

	return 1
}

// classifyCard() turns an single card string into its int representation.
// E.g. "A" -> 15.
func classifyCard(cardStr string, part PuzzlePart) int {
	// number-strings convert cleanly, and "A", "K", etc., get an error.
	result, err := strconv.Atoi(cardStr)

	// Card values change between parts 1 and 2.
	switch part {
	case Part1:
		if err != nil {
			switch cardStr {
			case "A":
				result = A1
			case "K":
				result = K1
			case "Q":
				result = Q1
			case "J":
				result = J1
			case "T":
				result = T1
			}
		}
	case Part2:
		if err != nil {
			switch cardStr {
			case "A":
				result = A2
			case "K":
				result = K2
			case "Q":
				result = Q2
			case "J":
				result = J2
			case "T":
				result = T2
			}
		}
	}

	return result
}

// classifyHand() returns the value of a hand type (e.g. 0 for High Card and 6
// for Five of a Kind).
func (h Hand) classifyHand() int {
	// Calculate how many times each card is seen - for determining hands.
	// Use h.wildCards if it's set (i.e. do Part2).
	cardMap := make(map[int]int)
	if len(h.wildCards) > 0 {
		for _, card := range h.wildCards {
			cardMap[card]++
		}
	} else {
		for _, card := range h.cards {
			cardMap[card]++
		}
	}

	// getHighestValue() gets the highest value in the cardMap.
	// Used for a secondary ranking to determine, for example, the difference
	// between Four of a Kind and a Full House, as both have two ranks of cards,
	// and the difference is the distribution of the ranks -- 4 of one rank vs 3.
	getHighestValue := func(h map[int]int) int {
		result := 0
		for _, v := range h {
			if v > result {
				result = v
			}
		}
		return result
	}

	switch len(cardMap) {
	case 1:
		return FiveOfAKind
	case 2:
		switch getHighestValue(cardMap) {
		case 4:
			return FourOfAKind
		default:
			return FullHouse
		}
	case 3:
		switch getHighestValue(cardMap) {
		case 3:
			return ThreeOfAKind
		default:
			return TwoPair
		}
	case 4:
		return OnePair
	case 5:
		return HighCard
	default:
		fmt.Println("Unrecognized hand")
		os.Exit(1)
	}

	if len(cardMap) == 3 {
		for _, v := range cardMap {
			if v == 3 {
				return ThreeOfAKind
			} else {
				return TwoPair
			}
		}
	}

	return 0
}

// isStronger() returns true if otherHand is stronger and false otherwise.
func (h Hand) isStrongerThan(otherHand Hand) bool {
	// When hand-type is identical, the hand with the highest first card is stronger.
	if otherHand.classifyHand() == h.classifyHand() {
		for i, card := range h.cards {
			if card == otherHand.cards[i] {
				continue
			}
			return card > otherHand.cards[i]
		}
	}

	// The cards weren't identical, so simply compare hand types.
	return h.classifyHand() > otherHand.classifyHand()
}

// getHands() converts lines of unparsed text into []Hand.
func getHands(lines []string, part PuzzlePart) Hands {
	regex := regexp.MustCompile(`(\d+|\w+) (\d+)`)
	result := []Hand{}

	for _, line := range lines {
		match := regex.FindStringSubmatch(line)
		handStr := match[1]
		bidStr := match[2]
		bid, err := strconv.Atoi(bidStr)
		if err != nil {
			panic(err)
		}

		cards := []int{}
		for _, cardStr := range handStr {
			cards = append(cards, classifyCard(string(cardStr), part))
		}

		result = append(result, newHand(cards, bid))
	}

	return result
}

// Sorting Hands by reverse rank.
func (h Hands) Len() int           { return len(h) }
func (h Hands) Less(i, j int) bool { return h[i].isStrongerThan(h[j]) }
func (h Hands) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Hands) sort() {
	sort.Sort(sort.Reverse(h))
}

// getTotalWinnings() returns the total winnings for a set of hands.
func (h Hands) getTotalWinnings() int {
	// First sort the hands before totaling.
	h.sort()
	total := 0
	for i, hand := range h {
		total += (i + 1) * hand.bid
	}

	return total
}

// Helper functions

func getLines(file io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// getCombinations() generates a [][]int with all possible combination of cards
// in an []int of len(slice).
func getCombinations(slice []int) [][]int {
	possibleCards := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	results := [][]int{{}}
	for range slice {
		newResult := [][]int{}
		for _, result := range results {
			for _, possibleCard := range possibleCards {
				resultWithNewCard := append(append([]int(nil), result...), possibleCard)
				newResult = append(newResult, resultWithNewCard)
			}
		}
		results = newResult
	}

	return results
}
