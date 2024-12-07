package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	INPUT_FILE = "input.txt"
	// INPUT_FILE = "testinput.txt"
	// INPUT_FILE = "manualinput.txt"
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

type Line struct {
	sum   int
	parts []int
}

func InputsFromLines(lines []string) []Line {
	inputs := []Line{}
	for _, line := range lines {
		sumStr := strings.Split(line, ":")[0]
		rest := strings.Split(line, ":")[1]
		sum, err := strconv.Atoi(sumStr)
		newParts := []int{}
		for _, part := range strings.Split(rest, " ") {
			if part == "" {
				continue
			}
			partInt, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			newParts = append(newParts, partInt)
		}
		if err != nil {
			panic(err)
		}
		l := Line{
			sum:   sum,
			parts: newParts,
		}

		inputs = append(inputs, l)
	}
	return inputs
}

type Operator string

func (o Operator) String() string {
	return string(o)
}

const (
	OP_ADD    = "+"
	OP_MUL    = "*"
	OP_CONCAT = "||"
)

func Operate(parts []int, slots []Operator) int {
	// fmt.Println("Operating with slots: ", slots)
	sum := 0
	for i, part := range parts {
		if i == 0 {
			sum = part
			continue
		}
		op := slots[i-1]
		switch op {
		case OP_ADD:
			sum += part
		case OP_MUL:
			sum *= part
		case OP_CONCAT:
			var err error
			concat := strconv.Itoa(sum) + strconv.Itoa(part)
			sum, err = strconv.Atoi(concat)
			if err != nil {
				panic(err)
			}
		}
	}
	return sum
}

func decimalToBase3(num int) string {
	if num == 0 {
		return "0"
	}
	result := ""
	for num > 0 {
		remainder := num % 3
		result = strconv.Itoa(remainder) + result
		num = num / 3
	}
	return result
}

const NUM_OP_TYPES = 3

func stringSolution(line Line, slots []Operator) string {
	solution := ""
	for i, part := range line.parts {
		if i == 0 {
			solution += strconv.Itoa(part)
			continue
		}
		solution += slots[i-1].String()
		solution += strconv.Itoa(part)
	}
	return solution
}

func validate(line Line) bool {
	fmt.Println("Validating: ", line)
	numOps := len(line.parts) - 1
	numVariations := int(math.Pow(NUM_OP_TYPES, float64(numOps)))
	slots := make([]Operator, numOps)
	// Default all slots to add
	for i := range slots {
		slots[i] = OP_ADD
	}
	fmt.Println("Checking variations: ", numVariations)
	for i := 0; i < numVariations; i++ {
		// Convert i to binary, then set slots to ADD if 0, MUL if 1
		value := decimalToBase3(i)
		// Pad value left with 0s
		for len(value) < numOps {
			value = "0" + value
		}
		for j := 0; j < numOps; j++ {
			if j >= len(slots) || j >= len(value) {
				slots[j] = OP_ADD
			} else if value[j] == '0' {
				slots[j] = OP_ADD
			} else if value[j] == '1' {
				slots[j] = OP_MUL
			} else if value[j] == '2' {
				slots[j] = OP_CONCAT
			}
		}

		sum := Operate(line.parts, slots)
		if sum == line.sum {
			fmt.Printf("Solution found: %s\n", stringSolution(line, slots))
			return true
		}
	}
	return false
}

func main() {
	lines := ReadFile(INPUT_FILE)

	inputs := InputsFromLines(lines)
	fmt.Printf("inputs: %v\n", inputs)

	validInputs := []Line{}
	for _, input := range inputs {
		if validate(input) {
			fmt.Printf("Valid!\n")
			validInputs = append(validInputs, input)
		}
	}
	fmt.Printf("valid inputs: %v\n", validInputs)
	sum := 0
	for _, input := range validInputs {
		sum += input.sum
	}
	fmt.Printf("Sum of valid inputs: %v\n", sum)
}
