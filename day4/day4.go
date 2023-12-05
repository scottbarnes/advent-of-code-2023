package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const filename = "day4_input.txt"

type FindType int

const (
	Points FindType = iota
	Copies
)

type Card struct {
	cardNo         int
	copies         int
	matches        int
	winningNumbers []int
	yourNumbers    []int
	deck           *CardDeck
}

type CardDeck struct {
	cards map[int]Card
}

// Global state. :(
// var cards = make(map[int]Card)

func NewCard(deck *CardDeck, line string, index int) Card {
	// Parse the line into a usable format.
	numberRegex := regexp.MustCompile(`\d+`)
	colonIndex := strings.Index(line, ":")
	lineWithoutCardNo := line[colonIndex+1:]
	winnersAndYourCards := strings.Split(lineWithoutCardNo, "|")
	winnerStrMatches := numberRegex.FindAllString(winnersAndYourCards[0], -1)
	yourStrMatches := numberRegex.FindAllString(winnersAndYourCards[1], -1)

	// Helper func to turn a []string{"1", "2", "3"} int []int{1, 2, 3}.
	makeNumbers := func(numStr []string) []int {
		nums := make([]int, len(numStr))
		for index, numStr := range numStr {
			parsedNum, _ := strconv.Atoi(numStr)
			nums[index] = parsedNum
		}
		return nums
	}

	winners := makeNumbers(winnerStrMatches)
	yourNumbers := makeNumbers(yourStrMatches)

	card := Card{
		cardNo:         index + 1,
		winningNumbers: winners,
		yourNumbers:    yourNumbers,
		copies:         1,
		deck:           deck,
	}
	card.matches = card.getMatches()

	return card
}

func NewCardDeck() *CardDeck {
	return &CardDeck{
		cards: make(map[int]Card),
	}
}

func (cd *CardDeck) addCard(card Card) {
	cd.cards[card.cardNo] = card
}

// contains returns true if a slice contains a given number.
func (c Card) contains(slice []int, number int) bool {
	for _, v := range slice {
		if v == number {
			return true
		}
	}
	return false
}

// getPoints returns the total points for a card (for part 1).
func (c Card) getPoints() int {
	if c.matches == 0 {
		return 0
	}

	if c.matches == 1 {
		return 1
	}

	result := 1
	if c.matches >= 1 {
		for i := 1; i < c.matches; i++ {
			result *= 2
		}
	}
	return result
}

// getMatches returns the number of matching numbers on a card.
func (c Card) getMatches() int {
	var matches int
	for _, v := range c.winningNumbers {
		if c.contains(c.yourNumbers, v) {
			matches += 1
		}
	}
	return matches
}

// For each copy of a card, add more cards based on the number of winning/matching numbers
// For example, if card 10 were to have 5 matching numbers, you would win one copy each
// of cards 11, 12, 13, 14, and 15, if you had TWO copies of card 10, you'd do this TWICE,
// ending up with more copies of 11-15.
func (c Card) addCopies() {
	start := c.cardNo + 1
	end := c.cardNo + c.matches
	for subsequentIndex := start; subsequentIndex <= end; subsequentIndex++ {
		// Don't add copies of cards beyond the limit.
		if subsequentIndex > len(c.deck.cards) {
			continue
		}
		card := c.deck.cards[subsequentIndex]
		card.copies += c.copies
		c.deck.cards[subsequentIndex] = card
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please provide the parameter 'points' or 'copies'.")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	switch os.Args[1] {
	case "points":
		result := run(file, Points)
		fmt.Println(result)
	case "copies":
		result := run(file, Copies)
		fmt.Println("copies:", result)
	}
}

// loadCardsFromFile reads file, creates one card per line, and adds them to deck.
func loadCardsFromFile(deck *CardDeck, file io.Reader) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for idx, line := range lines {
		card := NewCard(deck, line, idx)
		deck.addCard(card)
	}
}

func (cd *CardDeck) getTotalPoints() int {
	var total int
	for _, card := range cd.cards {
		total += card.getPoints()
	}
	return total
}

func (cd *CardDeck) getTotalcopies() int {
	var total int
	for i := 0; i < len(cd.cards); i++ {
		cd.cards[i].addCopies()
	}

	for _, card := range cd.cards {
		total += card.copies
	}

	return total
}

func run(file io.Reader, findType FindType) int {
	deck := NewCardDeck()
	loadCardsFromFile(deck, file)

	switch findType {
	case Points:
		return deck.getTotalPoints()
	case Copies:
		return deck.getTotalcopies()
	default:
		return 0
	}
}
