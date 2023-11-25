package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Constants for cell states
const (
	Dead       = 0
	Alive      = 1
	board_size = 160
)

func main() {
	// Define the 2D array
	rows := 8000
	cols := 16000

	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	// Initialize the array with some values (0 or 1)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrix[i][j] = i % 2
		}
	}

	rand.Seed(42)
	for i := range matrix {
		for j := range matrix[i] {
			// Randomly make some cell live
			randomNumber := rand.Intn(2)
			matrix[i][j] = randomNumber
		}
	}

	// Create a result matrix with the same dimensions
	resultMatrix := make([][]int, rows)
	for i := range resultMatrix {
		resultMatrix[i] = make([]int, cols)
	}

	// Use channels and wait groups for synchronization
	var wg sync.WaitGroup
	var muMatrix sync.Mutex
	var resultMutexRW sync.RWMutex

	// Game loop
	epochs := 10
	startTime := time.Now()
	elapsedEpoch := make([]int64, epochs)
	for i := 0; i < epochs; i++ {

		epochstartTime := time.Now()

		wg.Add((rows * cols) / (board_size * board_size))

		// Start producers (goroutines)
		for i := 0; i < rows; i += board_size {
			for j := 0; j < cols; j += board_size {
				wg.Add(1)
				go func(row, col int) {
					defer wg.Done()
					changeValue(&matrix, &muMatrix, &resultMatrix, &resultMutexRW,
						row, col, &wg)
				}(i, j)
			}
		}

		// Start a single consumer (goroutine)
		wg.Wait()

		matrix = resultMatrix

		elapsedEpoch[i] = time.Since(epochstartTime).Milliseconds()
		fmt.Println("epoca numero", i)

		//fmt.Print("\033[H\033[2J")
		//displayBoard(matrix)
		//time.Sleep(500 * time.Millisecond)
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Total Execution time: %s\n", elapsed)
	result := elapsed / time.Duration(epochs)
	fmt.Printf("Elapsed time divided by %d: %s\n", epochs, result)

	for i := 0; i < epochs; i++ {
		fmt.Println(elapsedEpoch[i])
	}
}

// Function to change the value of a matrix element
func changeValue(matrix *[][]int, muMatrix *sync.Mutex,
	resultMatrix *[][]int, muResult *sync.RWMutex,
	row, col int, wg *sync.WaitGroup) {

	defer muResult.RUnlock()
	defer wg.Done()

	muResult.RLock()

	muMatrix.Lock()
	for i := row; i < (board_size + row); i++ {
		for j := col; j < (board_size + col); j++ {
			neighbors := countNeighbors(matrix, i, j)
			result := Dead
			if (*matrix)[i][j] == Alive {
				// Cell is alive
				if neighbors < 2 || neighbors > 3 {
					// Loneliness or overcrowding, cell dies
					result = Dead
				} else {
					// Cell survives
					result = Alive
				}
			} else {
				// Cell is dead
				if neighbors == 3 {
					// Cell becomes alive
					result = Alive
				} else {
					// Cell remains dead
					result = Dead
				}
			}

			(*resultMatrix)[i][j] = result
		}
	}
	muMatrix.Unlock()
}

// Function to display a matrix
func displayMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}

// Function to count the number of live neighbors for a cell
func countNeighbors(board *[][]int, x, y int) int {
	rows := len((*board))
	cols := len((*board)[0])

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
				count += (*board)[nx][ny]
			}
		}
	}

	return count
}

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
