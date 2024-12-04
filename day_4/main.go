package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// INPUT_FILE  = "testinput.txt"
	INPUT_FILE  = "input.txt"
	TARGET_WORD = "XMAS"
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
	fmt.Printf("[%d][%d] => %s\n", i, j, proposedWord)
	isProposedWordPresent := strings.Index(TARGET_WORD, proposedWord)
	if isProposedWordPresent != 0 {
		// Can't find the proposed word
		return false
	}

	if proposedWord == TARGET_WORD {
		// We're done
		fmt.Printf("\t[%s] Match!\n", direction.String())
		return true
	}

	switch direction {
	case UP:
		if i < 0 {
			return false
		}
		return hunt(i-1, j, lines, proposedWord, direction)
	case UP_RIGHT:
		if i < 0 || j >= len(lines[i]) {
			return false
		}
		return hunt(i-1, j+1, lines, proposedWord, direction)

	case RIGHT:
		if j >= len(lines[i]) {
			return false
		}
		return hunt(i, j+1, lines, proposedWord, direction)
	case DOWN_RIGHT:
		if i >= len(lines) || j >= len(lines[i]) {
			return false
		}
		return hunt(i+1, j+1, lines, proposedWord, direction)
	case DOWN:
		if i >= len(lines) {
			return false
		}
		return hunt(i+1, j, lines, proposedWord, direction)
	case DOWN_LEFT:
		if i >= len(lines) || j < 0 {
			return false
		}
		return hunt(i+1, j-1, lines, proposedWord, direction)

	case LEFT:
		if j < 0 {
			return false
		}
		return hunt(i, j-1, lines, proposedWord, direction)
	case UP_LEFT:
		if i < 0 || j < 0 {
			return false
		}
		return hunt(i-1, j-1, lines, proposedWord, direction)

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

	// Up: i-1, j
	if hunt(i-1, j, lines, string(lines[i][j]), UP) {
		count++
	}
	// Up Right: i-1, j+1
	if hunt(i-1, j+1, lines, string(lines[i][j]), UP_RIGHT) {
		count++
	}
	// Right: i, j+1
	if hunt(i, j+1, lines, string(lines[i][j]), RIGHT) {
		count++
	}
	// Down Right: i+1, j+1
	if hunt(i+1, j+1, lines, string(lines[i][j]), DOWN_RIGHT) {
		count++
	}
	// Down: i+1, j
	if hunt(i+1, j, lines, string(lines[i][j]), DOWN) {
		count++
	}
	// Down Left: i+1, j-1
	if hunt(i+1, j-1, lines, string(lines[i][j]), DOWN_LEFT) {
		count++
	}
	// Left: i, j-1
	if hunt(i, j-1, lines, string(lines[i][j]), LEFT) {
		count++
	}
	// Up Left: i-1, j-1
	if hunt(i-1, j-1, lines, string(lines[i][j]), UP_LEFT) {
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
			if char == 'X' {
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
