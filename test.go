package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Constants for cell states
const (
	Dead  = 0
	Alive = 1
)

// Function to display the game board
func displayBoard(board [][]int) {
	for _, row := range board {
		for _, cell := range row {
			if cell == Alive {
				color.Set(color.BgBlack, color.Bold)
				fmt.Print("  ") // Print a filled square for alive cells
				color.Unset()
			} else {
				color.Set(color.BgHiWhite, color.Bold)
				fmt.Print("  ") // Print an empty square for dead cells
				color.Unset()
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Function to update the game board based on the rules of Conway's Game of Life
func updateBoard(board [][]int) [][]int {
	rows := len(board)
	cols := len(board[0])

	// Create a new board to store the next state
	newBoard := make([][]int, rows)
	for i := range newBoard {
		newBoard[i] = make([]int, cols)
	}

	// Apply the rules to update the new board
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			neighbors := countNeighbors(board, i, j)
			if board[i][j] == Alive {
				// Cell is alive
				if neighbors < 2 || neighbors > 3 {
					// Loneliness or overcrowding, cell dies
					newBoard[i][j] = Dead
				} else {
					// Cell survives
					newBoard[i][j] = Alive
				}
			} else {
				// Cell is dead
				if neighbors == 3 {
					// Cell becomes alive
					newBoard[i][j] = Alive
				} else {
					// Cell remains dead
					newBoard[i][j] = Dead
				}
			}
		}
	}

	return newBoard
}

// Function to count the number of live neighbors for a cell
func countNeighbors(board [][]int, x, y int) int {
	rows := len(board)
	cols := len(board[0])

	count := 0

	// Check the eight neighboring cells
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the cell itself
			if i == 0 && j == 0 {
				continue
			}

			// Calculate the neighbor's coordinates
			nx, ny := x+i, y+j

			// Check if the neighbor is within the bounds of the board
			if nx >= 0 && nx < rows && ny >= 0 && ny < cols {
				// Increment the count if the neighbor is alive
				count += board[nx][ny]
			}
		}
	}

	return count
}

func strToInt(str string) int {
	count, err := strconv.Atoi(str)
	if err != nil {
		// Handle the error, e.g., print an error message
		fmt.Println("Error converting count to an integer:", err)
	}
	return count
}

func RleToBoard(board [][]int, rle string, y int, x int) [][]int {
	//get the lines
	lines := strings.Split(rle, "\n")

	info := strings.Join(lines[1:], "")
	//get the row of the new board
	rows := strings.Split(info, "$")

	for i := 0; i < len(rows); i++ {
		countStr := ""
		xCounter := 0
		rows[i] += "$"
		for h, char := range rows[i] {
			if '0' <= char && char <= '9' {
				countStr += string(char)
			} else {
				count := 1
				if countStr != "" {
					count = strToInt(countStr)
				}
				if rows[i][h] == 'b' {
					for j := xCounter; j < (xCounter + count); j++ {
						board[y+i][x+j] = 0
					}
					xCounter = (xCounter + count)
				} else if (rows[i][h] == '$') && (count > 1) {
					y += count
				} else if rows[i][h] == 'o' {
					for j := xCounter; j < (xCounter + count); j++ {
						board[y+i][x+j] = 1
					}
					xCounter = (xCounter + count)
				}
				countStr = ""
			}
		}
	}
	return board
}

func main() {
	// Define the size of the game board
	rows, cols := 40, 60

	// Create the initial game board with random live cells
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}

	rle := `x = 47, y = 26, rule = B3/S23
bo$2bo$3o6$15bo$15b4o$16b4o10b2o$5b2o9bo2bo9bobo$5b2o9b4o8b3o8b2o$15b
4o8b3o9b2o$15bo12b3o$29bobo$30b2o7$45b2o$44b2o$46bo`
	RleToBoard(board, rle, 5, 5)

	displayBoard(board)

	// Game loop
	for {
		// Clear the console (for better visualization in the console)
		fmt.Print("\033[H\033[2J")

		// Display the current state of the game board
		displayBoard(board)

		// Update the game board based on the rules
		board = updateBoard(board)

		// Pause for a while to make it visually comprehensible
		time.Sleep(7000 * time.Millisecond)
	}

}
