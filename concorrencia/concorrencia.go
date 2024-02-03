package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"os"
	"strconv"
)

const (
	Dead       = 0
	Alive      = 1
	
)
var board_size int = 0

type Matrix [][]int

func main() {
	dim, _ := strconv.Atoi(os.Args[1])
	board_size, _ = strconv.Atoi(os.Args[2])
	epochs, _ := strconv.Atoi(os.Args[3])

	rows := dim
	cols := dim

	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	rand.Seed(42)
	for i := range matrix {
		for j := range matrix[i] {
			randomNumber := rand.Intn(2)
			matrix[i][j] = randomNumber
		}
	}

	resultMatrix := make([][]int, rows)
	for i := range resultMatrix {
		resultMatrix[i] = make([]int, cols)
	}

	matrixPool := &sync.Pool{
		New: func() interface{} {
			bufferMatrix := make(Matrix, board_size)
			for i := range bufferMatrix {
				bufferMatrix[i] = make([]int, board_size)
			}
			return bufferMatrix
		},
	}

	var wg sync.WaitGroup
	var muMatrix sync.Mutex
	var resultMutexRW sync.RWMutex

	elapsedEpoch := make([]int64, epochs)
	for i := 0; i < epochs; i++ {

		epochstartTime := time.Now()

		for i := 0; i < rows; i += board_size {
			for j := 0; j < cols; j += board_size {
				wg.Add(1)
				go func(row, col int) {
					defer wg.Done()
					changeValue(&matrix, &muMatrix, &resultMatrix, &resultMutexRW, matrixPool,
						row, col)
				}(i, j)
			}
		}
		wg.Wait()

		matrix = resultMatrix

		elapsedEpoch[i] = time.Since(epochstartTime).Microseconds()
	}

	file, _ := os.OpenFile("outputs/"+os.Args[1]+"_"+os.Args[2]+"_"+os.Args[3]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE,0222)
	for i := 0; i < epochs; i++ {
		fmt.Fprintf(file, "%d\n", elapsedEpoch[i])
	}
}
func changeValue(matrix *[][]int, muMatrix *sync.Mutex,
	resultMatrix *[][]int, muResult *sync.RWMutex, matrixPool *sync.Pool,
	row, col int) {

	defer muResult.RUnlock()

	bufferMatrix := matrixPool.Get().(Matrix)
	defer matrixPool.Put(bufferMatrix)

	muResult.RLock()

	for i := row; i < (board_size + row); i++ {
		for j := col; j < (board_size + col); j++ {
			neighbors := countNeighbors(matrix, i, j)
			result := Dead
			if (*matrix)[i][j] == Alive {
				if neighbors < 2 || neighbors > 3 {
					result = Dead
				} else {
					result = Alive
				}
			} else {
				if neighbors == 3 {
					result = Alive
				} else {
					result = Dead
				}
			}
			bufferMatrix[i-row][j-col] = result
		}
	}
	muMatrix.Lock()
	for i := row; i < (board_size + row); i++ {
		for j := col; j < (board_size + col); j++ {
			(*resultMatrix)[i][j] = bufferMatrix[i-row][j-col]
		}
	}
	muMatrix.Unlock()
}

func displayMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}

func countNeighbors(board *[][]int, x, y int) int {
	rows := len((*board))
	cols := len((*board)[0])

	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			nx, ny := x+i, y+j
			if nx >= 0 && nx < rows && ny >= 0 && ny < cols {
				count += (*board)[nx][ny]
			}
		}
	}
	return count
}
