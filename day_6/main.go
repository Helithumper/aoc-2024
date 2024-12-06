package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
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

type Object interface {
	String() string
}

type EmptyTile struct{}

func (EmptyTile) String() string {
	return "."
}

type Wall struct{}

func (Wall) String() string {
	return "#"
}

type Player struct {
	Direction Direction
	Position  [2]int
	Steps     map[[2]int]bool
}

func (p Player) String() string {
	switch p.Direction {
	case DIRECTION_UP:
		return "^"
	case DIRECTION_DOWN:
		return "v"
	case DIRECTION_LEFT:
		return "<"
	case DIRECTION_RIGHT:
		return ">"
	default:
		panic("Invalid direction")
	}
}

func (p Player) UniqueSteps() int {
	return len(p.Steps)
}

type Direction string

const (
	DIRECTION_UP    Direction = "^"
	DIRECTION_DOWN  Direction = "v"
	DIRECTION_LEFT  Direction = "<"
	DIRECTION_RIGHT Direction = ">"
)

type Board struct {
	originalLines        []string
	tiles                [][]Object
	player               Player
	blockerLoopPositions map[[2]int]bool
}

func (b Board) String() string {
	str := ""
	for i, row := range b.tiles {
		for j, tile := range row {
			isBlocker := b.blockerLoopPositions[[2]int{i, j}]
			if isBlocker {
				str += color.BlueString("B")
			} else if i == b.player.Position[0] && j == b.player.Position[1] {
				str += color.YellowString(string(b.player.Direction))
			} else if b.player.Steps[[2]int{i, j}] {
				str += color.RedString("x")
			} else {
				str += color.GreenString(tile.String())
			}
		}
		str += "\n"
	}
	return str
}

// Step moves the player one step in the current direction.
// returns true if the player has exited bounds
func (b *Board) Step(checkIsLooping bool) bool {
	var atPosition Object
	// Ensure current tile is marked as visited
	b.player.Steps[b.player.Position] = true

	isBlocker := false
	if checkIsLooping && b.willBlockerCauseLoop() {
		isBlocker = true
	}
	switch b.player.Direction {
	case DIRECTION_UP:
		if b.player.Position[0] == 0 {
			return true
		}
		atPosition = b.tiles[b.player.Position[0]-1][b.player.Position[1]]
		if atPosition.String() == "#" {
			// Turn 90 Degrees right
			b.player.Direction = DIRECTION_RIGHT
		} else {
			b.player.Position = [2]int{b.player.Position[0] - 1, b.player.Position[1]}
		}
	case DIRECTION_DOWN:
		if b.player.Position[0] == len(b.tiles)-1 {
			return true
		}
		atPosition = b.tiles[b.player.Position[0]+1][b.player.Position[1]]
		if atPosition.String() == "#" {
			// Turn 90 Degrees right
			b.player.Direction = DIRECTION_LEFT
		} else {
			b.player.Position = [2]int{b.player.Position[0] + 1, b.player.Position[1]}
		}
	case DIRECTION_LEFT:
		if b.player.Position[1] == 0 {
			return true
		}
		atPosition = b.tiles[b.player.Position[0]][b.player.Position[1]-1]
		if atPosition.String() == "#" {
			// Turn 90 Degrees right
			b.player.Direction = DIRECTION_UP
		} else {
			b.player.Position = [2]int{b.player.Position[0], b.player.Position[1] - 1}
		}
	case DIRECTION_RIGHT:
		if b.player.Position[1] == len(b.tiles[0])-1 {
			return true
		}
		atPosition = b.tiles[b.player.Position[0]][b.player.Position[1]+1]
		if atPosition.String() == "#" {
			// Turn 90 Degrees right
			b.player.Direction = DIRECTION_DOWN
		} else {
			b.player.Position = [2]int{b.player.Position[0], b.player.Position[1] + 1}
		}
	}
	if isBlocker {
		b.blockerLoopPositions[b.player.Position] = true
	}

	return false
}

func (b Board) willBlockerCauseLoop() bool {
	var blockerPosition [2]int
	switch b.player.Direction {
	case DIRECTION_UP:
		blockerPosition = [2]int{b.player.Position[0] - 1, b.player.Position[1]}
	case DIRECTION_DOWN:
		blockerPosition = [2]int{b.player.Position[0] + 1, b.player.Position[1]}
	case DIRECTION_LEFT:
		blockerPosition = [2]int{b.player.Position[0], b.player.Position[1] - 1}
	case DIRECTION_RIGHT:
		blockerPosition = [2]int{b.player.Position[0], b.player.Position[1] + 1}
	}

	if blockerPosition[0] < 0 || blockerPosition[1] < 0 || blockerPosition[0] >= len(b.tiles) || blockerPosition[1] >= len(b.tiles[0]) {
		return false
	}

	if b.tiles[blockerPosition[0]][blockerPosition[1]].String() == "#" {
		return false
	}

	// Create a new temporary board to check for loops
	tempBoard := Board{
		tiles: make([][]Object, len(b.tiles)),
		player: Player{
			Position:  b.player.Position,
			Direction: b.player.Direction,
			Steps:     make(map[[2]int]bool),
		},
		blockerLoopPositions: make(map[[2]int]bool), // not used in tempBoard
	}
	for i, row := range b.tiles {
		tempBoard.tiles[i] = make([]Object, len(row))
		for j, tile := range row {
			tempBoard.tiles[i][j] = tile
		}
	}

	tempBoard.tiles[blockerPosition[0]][blockerPosition[1]] = Wall{}
	isTempBoardLooping := tempBoard.isLooping()
	if !isTempBoardLooping {
		return false
	}

	// Double check that this placement didn't screw the board
	originalBoard := NewBoard(b.originalLines)
	originalBoard.tiles[blockerPosition[0]][blockerPosition[1]] = Wall{}
	if originalBoard.isLooping() {
		// original board also results in a loop, we're golden!
		return true
	}
	return false
}

type StepPosition struct {
	Position  [2]int
	Direction Direction
}

func (b Board) isLooping() bool {
	priorStepPositions := []StepPosition{}
	for {
		if b.Step(false) {
			return false
		}
		for _, p := range priorStepPositions {
			if p.Position[0] == b.player.Position[0] && p.Position[1] == b.player.Position[1] && p.Direction == b.player.Direction {
				return true
			}
		}
		priorStepPositions = append(priorStepPositions, StepPosition{b.player.Position, b.player.Direction})
	}
	// return false
}

func NewBoard(lines []string) Board {
	board := Board{
		originalLines: lines,
	}
	board.tiles = make([][]Object, len(lines))
	for i, line := range lines {
		board.tiles[i] = make([]Object, len(line))
		for j, char := range line {
			switch char {
			case '.':
				board.tiles[i][j] = EmptyTile{}
			case '#':
				board.tiles[i][j] = Wall{}
			case '^', 'v', '<', '>':
				p := Player{
					Position:  [2]int{i, j},
					Direction: Direction(string(char)),
					Steps:     make(map[[2]int]bool),
				}
				board.player = p
				board.tiles[i][j] = p
			}
		}
	}

	board.blockerLoopPositions = make(map[[2]int]bool)
	return board
}

type Game struct {
	board Board
}

func main() {
	lines := ReadFile(INPUT_FILE)
	fmt.Println(lines)
	b := NewBoard(lines)
	fmt.Printf("%s\n", b.String())
	// i := 1
	for !b.Step(true) {
	}
	fmt.Printf("[step %d]\n%s\n", b.player.UniqueSteps(), b.String())
	fmt.Printf("Solution: %d\n", b.player.UniqueSteps())
	fmt.Printf("Blocker Positions: %d\n", len(b.blockerLoopPositions))
}
