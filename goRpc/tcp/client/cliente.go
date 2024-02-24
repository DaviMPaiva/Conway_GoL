package main

import (
	"conway/goRpc/impl"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//args
	dim, _ := strconv.Atoi(os.Args[1])
	epochs, _ := strconv.Atoi(os.Args[2])
	print_result, _ := strconv.Atoi(os.Args[3])
	display, _ := strconv.Atoi(os.Args[4])
	file, _ := os.OpenFile("../../../outputs/GoRPC_"+os.Args[1]+"_"+os.Args[2]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)

	matrix := make([][]int, dim)
	//prepare matrix
	if display > 0 {
		matrix = plainTextReader(dim)
	} else {
		for i := range matrix {
			matrix[i] = make([]int, dim)
		}
		//rand init
		rand.Seed(int64(42))
		for i := range matrix {
			for j := range matrix[i] {
				randomNumber := rand.Intn(2)
				matrix[i][j] = randomNumber
			}
		}
	}
	// Conectar ao servidor RPC - host/porta
	matrix_aux := matrix
	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor", err)
		return
	}
	defer client.Close()

	for k := 0; k < int(epochs); k++ {
		//começa a contar o tempo
		start_time := time.Now()
		//request
		req := impl.Request{Matrix: matrix_aux, Dim: dim}
		rep := impl.Reply{}
		err = client.Call("ConwayGame.Initialize", req, &rep)
		if err != nil {
			fmt.Println("Erro na chamada remota:", err)
			return
		}
		// tempo decorrido
		elapsedTime := time.Since(start_time).Microseconds()
		//salva o tempo no arquivo
		fmt.Fprintf(file, "%d\n", elapsedTime)
		// atualiza a nova matrix
		matrix_aux = rep.Matrix_result
		// espera um tempo para printar, limpa o terminal e chama a funcao para printar
		if int(print_result) > 0 {
			time.Sleep(500)
			fmt.Println("\033[H\033[2J")
			displayBoard(matrix_aux)
		} else {
			fmt.Printf("pacote recebido numero %d\n", k)
		}
	}

}

const (
	Dead  = 0
	Alive = 1
)

func displayBoard(board [][]int) {
	for _, row := range board {
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

func plainTextReader(dim int) [][]int {
	// Input plaintext
	// Read file
	data, err := os.ReadFile("../../../Gosper-glider_gun.txt")
	if err != nil {
		fmt.Println(err)
	}

	// Convert to string
	plaintext := string(data)

	// Split into lines
	lines := strings.Split(plaintext, "\n")

	// Width is length of first line
	width := dim

	// Height is number of lines
	height := dim

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

	return matrix
}
