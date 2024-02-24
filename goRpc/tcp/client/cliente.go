package main

import (
	"conway/goRpc/impl"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func main() {
	//args
	dim, _ := strconv.Atoi(os.Args[1])
	epochs, _ := strconv.Atoi(os.Args[2])
	print_result, _ := strconv.Atoi(os.Args[3])
	file, _ := os.OpenFile("../../../outputs/GoRPC_"+os.Args[1]+"_"+os.Args[2]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)
	//prepare matrix
	matrix := make([][]int, dim)
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
			time.Sleep(time.Second * 2)
			fmt.Println("\033[H\033[2J")
			displayBoard(matrix_aux)
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
