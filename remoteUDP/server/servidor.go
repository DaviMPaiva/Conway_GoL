package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Data struct {
	Size   int     `json:"size"`
	Matrix [][]int `json:"matrix"`
}

var board_size = 10

func handleConnection(conn *net.UDPConn) {

	buffer := make([]byte, 65535)

	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	var receivedData Data

	err = json.Unmarshal(buffer[:n], &receivedData)

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
		conn.WriteToUDP(byted_result, addr)
	}
}

func convert_to_byte_arr(raw_matrix [][]int, dim int) []byte {
	byte_matrix := make([][]byte, dim*dim)
	for i := range raw_matrix {
		for j := range raw_matrix[0] {
			byte_matrix = append(byte_matrix, []byte(strconv.Itoa(raw_matrix[i][j])))
		}
	}
	byte_arr := bytes.Join(byte_matrix, nil)
	return byte_arr
}

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":5151")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Server listening on :5151")

	for {
		handleConnection(conn)
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
