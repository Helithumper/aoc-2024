package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// INPUT_FILE = "testinput.txt"
	INPUT_FILE  = "input.txt"
	TARGET_WORD = "MS" // We want to find SAM/MAS, however since we'll start at A we represent this as "MS"
)

func ReadFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

type Direction int

const (
	UP Direction = iota
	UP_RIGHT
	RIGHT
	DOWN_RIGHT
	DOWN
	DOWN_LEFT
	LEFT
	UP_LEFT
)

func (d Direction) String() string {
	return []string{"UP", "UP_RIGHT", "RIGHT", "DOWN_RIGHT", "DOWN", "DOWN_LEFT", "LEFT", "UP_LEFT"}[d]
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func hunt(i int, j int, lines []string, soFar string, direction Direction) bool {
	// Is this a valid position?
	if i < 0 || j < 0 || i >= len(lines) || j >= len(lines[i]) {
		return false
	}
	// Get the next character that we are looking for
	// Then check that the cureent character is the next character
	// then continue moving in the correct direction
	currentChar := string(lines[i][j])
	proposedWord := soFar + currentChar
	fmt.Printf("[%d][%d]<%s> => %s\n", i, j, direction.String(), proposedWord)
	isProposedWordPresentFwds := strings.Index(TARGET_WORD, proposedWord)
	isProposedWordPresentBkwds := strings.Index(TARGET_WORD, Reverse(proposedWord))
	fmt.Printf("\t[%s] Fwds: %d Bkwds: %d\n", direction.String(), isProposedWordPresentFwds, isProposedWordPresentBkwds)
	if isProposedWordPresentFwds == -1 && isProposedWordPresentBkwds == -1 {
		// Can't find the proposed word
		fmt.Printf("\t[%s] No Match\n", direction.String())
		return false
	}

	if proposedWord == TARGET_WORD || Reverse(proposedWord) == TARGET_WORD {
		// We're done
		fmt.Printf("\t[%s] Match!\n", direction.String())
		return true
	}

	// Since we want to explore the opposite side from where we are looking, these will be flipped
	// UP_RIGHT <-> DOWN_LEFT
	// DOWN_RIGHT <-> UP_LEFT
	// We also want to skip the original spot so we don't see the A
	switch direction {
	case DOWN_LEFT: // Was UP_RIGHT, now DOWN_LEFT
		return hunt(i-2, j+2, lines, proposedWord, direction)
	case UP_LEFT: // Was DOWN_RIGHT, now UP_LEFT
		return hunt(i+2, j+2, lines, proposedWord, direction)
	case UP_RIGHT: // Was DOWN_LEFT, now UP_RIGHT
		return hunt(i+2, j-2, lines, proposedWord, direction)
	case DOWN_RIGHT: // Was UP_LEFT, now DOWN_RIGHT
		return hunt(i-2, j-2, lines, proposedWord, direction)

	default:
		panic("invalid direction")
	}
}

func Hunt(i int, j int, lines []string) int {
	// Beginning at i,j we will move in a clockwise
	// direction, building our resulting string as we go
	// We will return the count of XMAS we find
	soFar := []string{}
	soFar = append(soFar, string(lines[i][j]))
	count := 0
	fmt.Printf("Hunting at [%d][%d]\n", i, j)

	// We have to see both to see an X
	if hunt(i-1, j+1, lines, "", UP_RIGHT) && hunt(i+1, j+1, lines, "", DOWN_RIGHT) {
		count++
	}

	return count
}

func TraverseGrid(lines []string) int {
	// Starting at the top left, we work our way right then down
	// If we find an X, we begin a HUNT!
	count := 0
	for i, line := range lines {
		for j, char := range line {
			if char == 'A' {
				count += Hunt(i, j, lines)
			}
		}
	}

	return count
}

func main() {
	lines := ReadFile(INPUT_FILE)

	count := TraverseGrid(lines)
	fmt.Println(count)
}
