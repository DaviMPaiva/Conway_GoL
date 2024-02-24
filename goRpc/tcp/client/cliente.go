package main

import (
	"conway/goRpc/impl"
	"fmt"
	"math/rand"
	"net/rpc"
	"time"
)

const dim = 40

func main() {
	matrix := make([][]int, dim)
	for i := range matrix {
		matrix[i] = make([]int, dim)
	}

	rand.Seed(int64(42))
	for i := range matrix {
		for j := range matrix[i] {
			randomNumber := rand.Intn(2)
			matrix[i][j] = randomNumber
		}
	}
	Cliente(matrix, dim)
}

func Cliente(matrix [][]int, dim int) {
	// 1: Conectar ao servidor RPC - host/porta
	matrix_aux := matrix
	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor", err)
		return
	}
	defer client.Close()
	for k := 0; k < 30; k++ {
		req := impl.Request{Matrix: matrix_aux, Dim: dim}
		//a funçao reply diz o que vai ser retornado da função add da calculadora
		rep := impl.Reply{}
		err = client.Call("ConwayGame.Initialize", req, &rep)
		if err != nil {
			fmt.Println("Erro na chamada remota:", err)
			return
		}

		// 3: Imprimir o resultado
		matrix_aux = rep.Matrix_result

		time.Sleep(time.Second * 2)
		fmt.Println("\033[H\033[2J")
		displayBoard(matrix_aux)

		//fmt.Printf("Request n° %d\n\n", k+1)
		//for i := 0; i < len(matrix_aux); i++ {
		//	for j := 0; j < len(matrix_aux[0]); j++ {
		//		fmt.Printf("%d ", matrix_aux[i][j])
		//	}
		//	fmt.Printf("\n")
		//}
	}

	// 2: Invocar a operação remota

}

const (
	Dead  = 0
	Alive = 1
)

func displayBoard(board [][]int) {
	// Loop through rows
	for _, row := range board {
		// Loop through columns
		for _, val := range row {
			// Print black or black spot based on value
			if val == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print("██")
			}
		}
		// Newline after each row
		fmt.Println()
	}
}
