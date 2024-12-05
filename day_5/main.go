package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	INPUT_FILE = "input.txt"
	// INPUT_FILE = "testinput.txt"
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

// generateRules takes the input and generates the
// pre and post rules for the given set.
// PRe rules map a node to the node that has to come AFTER it
// Post rules map a node to the nodes that have to come BEFORE it
func generateRules(lines []string) (map[int][]int, map[int][]int) {
	prePostRules := make(map[int][]int)
	postPreRules := make(map[int][]int)
	for _, line := range lines {
		if line == "" {
			break
		}
		// Split the line into the pre and
		// post nodes
		split := strings.Split(line, "|")
		pre, err := strconv.Atoi(split[0])
		if err != nil {
			panic("Error converting pre to int")
		}
		post, err := strconv.Atoi(split[1])
		if err != nil {
			panic("Error converting post to int")
		}

		if _, ok := prePostRules[pre]; ok {
			prePostRules[pre] = append(prePostRules[pre], post)
		} else {
			prePostRules[pre] = []int{post}
		}

		if _, ok := postPreRules[post]; ok {
			postPreRules[post] = append(postPreRules[post], pre)
		} else {
			postPreRules[post] = []int{pre}
		}
	}

	// Add the pre and post nodes to the rules
	return prePostRules, postPreRules
}

func toIntSlice(s string) []int {
	ints := []int{}
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic("Error converting string to int")
		}
		ints = append(ints, i)
	}
	return ints
}

func swap(line []int, a int, b int) []int {
	tmp := line[a]
	line[a] = line[b]
	line[b] = tmp
	return line
}

// Returns fixed line, and whether it was fixed or not
func fixedLine(line []int, prePostRules map[int][]int, postPreRules map[int][]int, fixed bool) ([]int, bool) {
	for n, node := range line {
		for c, cursor := range line {
			// For each pairup, check if the pre and post rules are
			// satisfied

			// Equality, we ignore as this is the same node
			if n == c {
				continue
			}

			// BEFORE current node, left of it
			// We want to ensure no postPreRules appear in the preceeding nodes
			if c < n {
				// NOTE flipped this and I ahve no idea why it worked
				// The node we are looking for is BEFORE (pre) the current node
				// we want to ensure that this node is not in the Posts list
				for _, preceding := range prePostRules[node] {
					if preceding == cursor {
						fmt.Printf("[n=%d][c=%d] Rule Violated: %d must come after %d\n", n, c, node, cursor)
						line = swap(line, n, c)
						return fixedLine(line, prePostRules, postPreRules, true)
					}
				}
			}

			if n < c {
				// The node we are looking for is AFTER (post) the current node
				// we want to ensure that this node is not in the Pres list
				for _, following := range postPreRules[node] {
					if following == cursor {
						fmt.Printf("[n=%d][c=%d] Rule Violated: n=%d must come before c=%d\n", n, c, node, cursor)
						line = swap(line, n, c)
						return fixedLine(line, prePostRules, postPreRules, true)
					}
				}
			}
		}
	}
	return line, fixed
}

// Find invalid orders and fix them, then return the fixed versions
func filterInValidOrders(lines []string, prePostRules map[int][]int, postPreRules map[int][]int) [][]int {
	validOrders := [][]int{}
	for _, line := range lines {
		if !strings.Contains(line, ",") {
			continue
		}

		// Convert the strint to a slice of ints
		ints := toIntSlice(line)
		fmt.Printf("Evaluating: %v\n", ints)
		fixed, isFixed := fixedLine(ints, prePostRules, postPreRules, false)
		if isFixed {
			validOrders = append(validOrders, fixed)
		}
	}
	return validOrders
}

func main() {
	lines := ReadFile(INPUT_FILE)
	pres, posts := generateRules(lines)
	fmt.Printf("Pres: %v\n", pres)
	fmt.Printf("Posts: %v\n", posts)

	valid := filterInValidOrders(lines, pres, posts)
	fmt.Println(valid)
	sum := 0
	for _, v := range valid {
		if len(v)%2 == 0 {
			panic(fmt.Sprintf("Even length result has no middle: %v", v))
		}
		fmt.Printf("VALID: %v\n", v)
		fmt.Printf("\tMIDDLE: %d\n", v[len(v)/2])
		sum += v[len(v)/2]
	}
	fmt.Printf("Sum of middles: %d\n", sum)
}
