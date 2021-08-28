package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// When our program starts let's build the universe. We need a board.
	board := newBoard()

	// We'll loop forever. In this version there is no way to end the game.
	for {
		// Everytime we iterate we will render the board.
		board.Render()
		// After we render the board, let's ask the player what move they
		// wish to make.
		err := board.Turn()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}
}

// cellStatus is a custom type to help us avoid programming errors. We will
// declare a few different known states for a cell below and reference them
// when we update the board and render it.
type cellStatus int

const (
	Available cellStatus = iota
	Player1Occupied
	Player2Occupied
	Block
)

// Cell represents a single point or position on the board. Every cell
// can maintain a state. In this iteration of the program we can only
// track the Status of a cell. Perhaps other states will be recorded in
// the future! One never knows after all what the future has in store for us!
type Cell struct {
	Status cellStatus
}

// Render converts the internal state of a cell into a UI representation.
func (c *Cell) Render() string {
	switch c.Status {
	case Available:
		return ""
	case Player1Occupied:
		return "X"
	case Player2Occupied:
		return "O"
	case Block:
		return "~"
	default:
		panic("unknown cell status")
	}
}

// We actually have a gridSize of `9` but we extend to allow the header and
// index to be printed onto the game board.
const gridSize = 9

func newBoard() Board {
	rows := make([][]Cell, 0, gridSize)
	for i := 0; i < gridSize; i++ {
		rows = append(rows, newRow())
	}
	return Board{
		Rows:          rows,
		playerOneTurn: true,
	}
}

func newRow() []Cell {
	columns := make([]Cell, 0, gridSize)
	for i := 0; i < gridSize; i++ {
		columns = append(columns, newCell())
	}
	return columns
}

func newCell() Cell {
	return Cell{
		Status: Available,
	}
}

// Board represents our main state.
type Board struct {
	// Rows manages the state of each individual cell.
	Rows          [][]Cell
	playerOneTurn bool

	playerOneRow int
	playerOneCol int
	playerTwoRow int
	playerTwoCol int
}

// Turn allows a player to take a turn.
func (b *Board) Turn() error {
	input := b.captureInput()
	row, column, err := b.validateInput(strings.TrimSpace(input))
	if err != nil {
		return err
	}
	b.move(row, column)
	b.playerOneTurn = !b.playerOneTurn // update state for the next player
	return nil
}

// captureInput prompts the user to make a move.
func (b *Board) captureInput() string {
	reader := bufio.NewReader(os.Stdin)
	playerTurn := "Player One"
	if !b.playerOneTurn {
		playerTurn = "Player Two"
	}
	fmt.Printf("[%s] Where would you like to move to?: ", playerTurn)
	s, _ := reader.ReadString('\n')
	return s
}

// validateInput confirms the user provided input looks ok.
func (b *Board) validateInput(text string) (int, int, error) {
	if len(text) != 2 {
		return -1, -1, fmt.Errorf("... ummm.... that is not a valid position. Try something like A7")
	}
	letter := text[0]
	row, err := strconv.Atoi(string(text[1]))
	if err != nil {
		return -1, -1, fmt.Errorf("... ummm.... are you taking this serious? Enter a position like A7", row)
	}
	if row > gridSize || row < 1 {
		return -1, -1, fmt.Errorf("... ummm.... row %d does not exist on the board", row)
	}
	column, ok := columnMapInverted[string(letter)]
	if !ok {
		return -1, -1, fmt.Errorf("... ummm.... column %s does not exist on the board", string(letter))
	}
	return row, column, nil
}

// move updates the internal state of the board based on the player's
// validated move.
func (b *Board) move(moveToRow, moveToColumn int) {
	if b.playerOneTurn {
		if b.playerOneRow > 0 {
			b.Rows[b.playerOneRow][b.playerOneCol] = Cell{Status: Block}
		}
		b.Rows[moveToRow][moveToColumn] = Cell{Status: Player1Occupied}
		b.playerOneRow = moveToRow
		b.playerOneCol = moveToColumn
		return
	}

	if b.playerTwoRow > 0 {
		b.Rows[b.playerTwoRow][b.playerTwoCol] = Cell{Status: Block}
	}
	b.Rows[moveToRow][moveToColumn] = Cell{Status: Player2Occupied}
	b.playerTwoRow = moveToRow
	b.playerTwoCol = moveToColumn
	return
}

// Render prints the entire board to stdout.
func (b *Board) Render() {
	table := tablewriter.NewWriter(os.Stdout)

	for i, row := range b.Rows {
		renderedRow := make([]string, 0, gridSize)
		for j, cell := range row {
			// Render the top-left cell
			if i == 0 && j == 0 {
				renderedRow = append(renderedRow, "")
				continue
			}
			// Render the row number
			if i == 0 {
				renderedRow = append(renderedRow, columnMap[j-1])
				continue
			}
			// Render the column letter
			if j == 0 {
				renderedRow = append(renderedRow, strconv.Itoa(i))
				continue
			}
			// Render an interior column
			renderedRow = append(renderedRow, cell.Render())
		}
		table.Append(renderedRow)
	}
	table.Render()
}

var columnMap = make(map[int]string)
var columnMapInverted = make(map[string]int)

func init() {
	for i := 0; i < gridSize; i++ {
		letter := string(rune(65 + i))
		columnMap[i] = letter
		columnMapInverted[letter] = i + 1
	}
}
