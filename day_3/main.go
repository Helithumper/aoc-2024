package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/fatih/color"
)

const INPUT_FILE = "testinput.txt"

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

func ExtractCleanText(text string) []string {
	re := regexp.MustCompile(`(mul\(\d{1,3},\d{1,3}\))|(do\(\))|(don\'t\(\))`)
	matches := re.FindAllString(text, -1)
	return matches
}

func ExecuteMulOp(op string) int {
	re := regexp.MustCompile(`\d{1,3}`)
	matches := re.FindAllString(op, -1)
	if len(matches) != 2 {
		log.Fatal("Invalid mul operation")
	}
	a, err := strconv.Atoi(matches[0])
	b, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal("Invalid number")
	}
	return a * b
}

func isDoFlag(op string) bool {
	return op == "do()"
}

func isDontFlag(op string) bool {
	return op == "don't()"
}

func main() {
	lines := ReadFile(INPUT_FILE)

	result := 0
	doFlag := true
	for _, line := range lines {
		cleanOps := ExtractCleanText(line)
		fmt.Println(cleanOps)

		for _, op := range cleanOps {
			if isDoFlag(op) {
				fmt.Printf("%s\n", color.GreenString("DO!"))
				doFlag = true
				continue
			} else if isDontFlag(op) {
				fmt.Printf("%s\n", color.RedString("DONT!"))
				doFlag = false
				continue
			}
			if doFlag {
				fmt.Printf("\tExecuting: %s\n", op)
				result += ExecuteMulOp(op)
			} else {
				fmt.Printf("\tSkipping: %s\n", op)
			}
		}
	}

	fmt.Printf("Result: %d\n", result)
}
