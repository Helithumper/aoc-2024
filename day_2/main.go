package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const INPUT_FILE = "input.txt"

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

func LinesToInts(lines []string) [][]int {
	result := [][]int{}
	for _, line := range lines {
		items := strings.Split(line, " ")
		linenums := []int{}
		for _, num := range items {
			v, err := strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
			linenums = append(linenums, v)
		}

		result = append(result, linenums)
	}
	return result
}

func eliminateIndex(line []int, idx int) []int {
	ret := make([]int, 0)
	ret = append(ret, line[:idx]...)
	res := append(ret, line[idx+1:]...)
	fmt.Printf("%v -> %v | Eliminated %v at index %v\n", line, res, line[idx], idx)
	return res
	/*
		   BUG: Some issue with regards to copying data. Multiple times seen where line gets overridden by append
			 https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			cpy := make([]int, len(line))
			cnt := copy(cpy, line)
			if cnt != len(line) {
				panic("Copy Failed")
			}
			res := append(cpy[:idx], cpy[idx+1:]...)
			fmt.Printf("%v -> %v | Eliminated %v at index %v\n", cpy, res, cpy[idx], idx)
			return res */
}

func DefuseCheck(line []int, idx int) bool {
	for i := 0; i < len(line); i++ {
		l := eliminateIndex(line, i)
		if IsSafe(l, false) {
			fmt.Println("DEFUSED")
			return true
		}
	}
	return false
}

// True if safe, False otherwise
// DEFUSE: allow one unsafe situation
func IsSafe(line []int, defuse bool) bool {
	increasing := line[0] < line[1]
	for i := 0; i < len(line)-1; i++ {
		// If actually decreasing or stable
		if line[i] >= line[i+1] && increasing {
			fmt.Printf("FAIL: %v >= %v and pattern is INCREASING\n", line[i], line[i+1])
			if defuse && DefuseCheck(line, i) {
				return true
			} else {
				return false
			}
		} else if line[i] <= line[i+1] && !increasing {
			fmt.Printf("FAIL: %v <= %v and pattern is DECREASING\n", line[i], line[i+1])
			if defuse && DefuseCheck(line, i) {
				return true
			} else {
				return false
			}
		}

		// Check difference
		diff := abs(line[i+1] - line[i])
		if diff < 1 || diff > 3 {
			fmt.Printf("FAIL: %v - %v = %v\n", line[i+1], line[i], diff)
			if defuse && DefuseCheck(line, i) {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func SafeCount(lines [][]int, defuse bool) int {
	cnt := 0
	for _, line := range lines {
		fmt.Printf("%s: %v\n", color.BlueString("Checking Line"), line)
		safe := IsSafe(line, defuse)
		res := ""
		if !safe {
			res = color.RedString("UNSAFE")
			fmt.Printf("Line %v is safe: %s\n", line, res)
		}
		if safe {
			cnt++
		}
	}

	return cnt
}

func main() {
	lines := ReadFile(INPUT_FILE)

	int_lines := LinesToInts(lines)

	safeCount := SafeCount(int_lines, true)
	fmt.Println("Safe Count: ", safeCount)
}
