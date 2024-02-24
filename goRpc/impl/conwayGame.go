package impl

import (
	"sync"
)

const (
	Dead  = 0
	Alive = 1
)

type Matrix [][]int

type ConwayGame struct{}

type Request struct {
	Matrix [][]int
	Dim    int
}

type Reply struct {
	Matrix_result [][]int
}

func (c *ConwayGame) Initialize(req Request, res *Reply) error {

	dim := req.Dim
	boardSize := 10
	rows := dim
	cols := dim

	matrix := req.Matrix

	resultMatrix := make(Matrix, rows)
	for i := range resultMatrix {
		resultMatrix[i] = make([]int, cols)
	}

	matrixPool := &sync.Pool{
		New: func() interface{} {
			bufferMatrix := make(Matrix, boardSize)
			for i := range bufferMatrix {
				bufferMatrix[i] = make([]int, boardSize)
			}
			return bufferMatrix
		},
	}

	var wg sync.WaitGroup
	var muMatrix sync.Mutex
	var resultMutexRW sync.RWMutex

	for i := 0; i < rows; i += boardSize {
		for j := 0; j < cols; j += boardSize {
			wg.Add(1)
			go func(row, col int) {
				defer wg.Done()
				changeValue(matrix, &muMatrix, &resultMatrix, &resultMutexRW, matrixPool,
					row, col, boardSize)
			}(i, j)
		}
	}
	wg.Wait()

	matrix = resultMatrix
	res.Matrix_result = matrix
	return nil
}

func changeValue(matrix [][]int, muMatrix *sync.Mutex,
	resultMatrix *Matrix, muResult *sync.RWMutex, matrixPool *sync.Pool,
	row, col int, boardSize int) {

	defer muResult.RUnlock()

	bufferMatrix := matrixPool.Get().(Matrix)
	defer matrixPool.Put(bufferMatrix)

	muResult.RLock()

	for i := row; i < (boardSize + row); i++ {
		for j := col; j < (boardSize + col); j++ {
			neighbors := countNeighbors(matrix, i, j)
			result := Dead
			if (matrix)[i][j] == Alive {
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
	for i := row; i < (boardSize + row); i++ {
		for j := col; j < (boardSize + col); j++ {
			(*resultMatrix)[i][j] = bufferMatrix[i-row][j-col]
		}
	}
	muMatrix.Unlock()
}

func countNeighbors(board [][]int, x, y int) int {
	rows := len((board))
	cols := len((board)[0])

	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			nx, ny := x+i, y+j
			if nx >= 0 && nx < rows && ny >= 0 && ny < cols {
				count += (board)[nx][ny]
			}
		}
	}
	return count
}
