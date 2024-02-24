package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Cell represents a cell in the game of life
type Cell struct {
	x, y int
}

func main() {
	// Provide the path to your RLE file
	filePath := "file.rle"

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Initialize variables for x, y, and cells
	var x, y int
	var size_x, size_y int
	var cells []Cell

	// Parse the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse x, y, and rule information
		if strings.HasPrefix(line, "x") {
			fmt.Sscanf(line, "x = %d, y = %d", &x, &y)
			continue
		}
		fmt.Printf("X = %d and Y = %d \n", x, y)

		// Parse cell information
		for _, char := range line {
			switch char {
			case 'b', 'o':
				cells = append(cells, Cell{x, y})
				x++
			case '$':
				x = 0
				y++
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				count := int(char - '0')
				for i := 0; i < count; i++ {
					cells = append(cells, Cell{x, y})
					x++
				}
			}
		}
	}
	fmt.Println(cells)
	// Create and initialize the 2D grid
	grid := make([][]bool, y+1)
	for i := range grid {
		grid[i] = make([]bool, x+1)
	}

	// Set cells in the grid to true based on the parsed cell coordinates
	for _, cell := range cells {
		grid[cell.y][cell.x] = true
	}

	// Print the game board
	printGrid(grid)
}

// Function to print the game board
func printGrid(grid [][]bool) {
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				fmt.Print("o")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
