package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Input plaintext
	data, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println(err)
	}

	// Convert to string
	plaintext := string(data)

	// Split into lines
	lines := strings.Split(plaintext, "\n")

	for _, line := range lines {
		fmt.Println(line)
	}

	// Width is length of first line
	width := len(lines[0]) + 3

	// Height is number of lines
	height := len(lines) + 3

	// Create matrix
	matrix := make([][]int, height)

	for i := range matrix {
		matrix[i] = make([]int, width)
	}

	// Parse each line into matrix
	for y, line := range lines {
		for x, c := range line {
			if c == 'O' {
				matrix[y][x] = 1
			}
		}
	}

	for _, row := range matrix {
		for _, val := range row {
			if val == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print("██")
			}
		}
		fmt.Println()
	}
}
