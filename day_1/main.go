package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	INPUT_FILE = "input.txt"
	LEFT_SIDE  = 0
	RIGHT_SIDE = 1
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

func CreatePairs(lines []string) [][]int {
	result := [][]int{{}, {}}
	for _, line := range lines {
		line := strings.Split(line, "   ")
		left, err := strconv.Atoi(line[LEFT_SIDE])
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(line[RIGHT_SIDE])
		if err != nil {
			panic(err)
		}
		result[LEFT_SIDE] = append(result[LEFT_SIDE], left)
		result[RIGHT_SIDE] = append(result[RIGHT_SIDE], right)
	}
	return result
}

func SortPairs(pairs [][]int) [][]int {
	slices.SortFunc(pairs[LEFT_SIDE], func(a int, b int) int { return a - b })
	slices.SortFunc(pairs[RIGHT_SIDE], func(a int, b int) int { return a - b })
	return pairs
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func Distances(pairs [][]int) []int {
	distances := []int{}
	for i := 0; i < len(pairs[LEFT_SIDE]); i++ {
		distances = append(distances, abs(pairs[RIGHT_SIDE][i]-pairs[LEFT_SIDE][i]))
	}
	return distances
}

func ComputeSum(distances []int) int {
	sum := 0
	for _, distance := range distances {
		sum += distance
	}
	return sum
}

func ComputeSimilarityScore(pairs [][]int) int {
	left_cnts := make(map[int]int)
	right_cnts := make(map[int]int)

	// Calculate scores based on the values in the right list
	// ASSUMPTION: Both slices are same size
	if len(pairs[LEFT_SIDE]) != len(pairs[RIGHT_SIDE]) {
		panic("Left and right slices are not the same size")
	}

	for i := 0; i < len(pairs[LEFT_SIDE]); i++ {
		l := pairs[LEFT_SIDE][i]
		r := pairs[RIGHT_SIDE][i]

		if val, ok := left_cnts[l]; ok {
			left_cnts[l] = val + 1
		} else {
			left_cnts[l] = 1
		}

		if val, ok := right_cnts[r]; ok {
			right_cnts[r] = val + 1
		} else {
			right_cnts[r] = 1
		}
	}

	// Score = SUM(LEFT_CNT * VAL * RIGHT_CNT)
	score := 0
	for key := range left_cnts {
		score += left_cnts[key] * key * right_cnts[key]
	}

	return score
}

func main() {
	// Read in file lines
	lines := ReadFile(INPUT_FILE)

	pairs := CreatePairs(lines)

	sortedPairs := SortPairs(pairs)

	distances := Distances(sortedPairs)

	sumOfDistances := ComputeSum(distances)

	fmt.Printf("Part 1 solution: Sum of distances: %d\n", sumOfDistances)

	similarityScore := ComputeSimilarityScore(sortedPairs)

	fmt.Printf("Part 2 solution: similarityScore: %d\n", similarityScore)
}
