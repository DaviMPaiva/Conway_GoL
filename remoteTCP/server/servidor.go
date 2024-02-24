package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Data struct {
	Size   int     `json:"size"`
	Matrix [][]int `json:"matrix"`
}

var board_size = 10

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)

		n, _ := conn.Read(buffer)

		var receivedData Data

		err := json.Unmarshal(buffer[:n], &receivedData)

		if err != nil {

			fmt.Println("Error unmarshaling JSON:", err)

		} else {

			//fmt.Println("Unmarshaled data:", receivedData)

			raw_result := conway_game(receivedData.Size, board_size, &receivedData.Matrix)
			data := Data{
				Size:   receivedData.Size,
				Matrix: raw_result,
			}
			byted_result, _ := json.Marshal(data)
			conn.Write(byted_result)
		}
	}
}

func main() {

	r, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("Error resolving:", err)
		return
	}
	//cria um listener tcp
	ln, err := net.ListenTCP("tcp", r)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

const (
	Dead  = 0
	Alive = 1
)

type Matrix [][]int

func conway_game(dim int, board_size int, matrix *[][]int) [][]int {

	rows := dim
	cols := dim

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

	for i := 0; i < rows; i += board_size {
		for j := 0; j < cols; j += board_size {
			wg.Add(1)
			go func(row, col int) {
				defer wg.Done()
				changeValue(matrix, &muMatrix, &resultMatrix, &resultMutexRW, matrixPool,
					row, col, board_size)
			}(i, j)
		}
	}
	wg.Wait()

	return resultMatrix
}
func changeValue(matrix *[][]int, muMatrix *sync.Mutex,
	resultMatrix *[][]int, muResult *sync.RWMutex, matrixPool *sync.Pool,
	row, col int, board_size int) {

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
