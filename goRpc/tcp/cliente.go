package main

import (
	"fmt"
	"impl"
	"math/rand"
	"net/rpc"
	"time"
)

const dim = 100

func main() {
	matrix := make([][]int, dim)
	for i := range matrix {
		matrix[i] = make([]int, dim)
	}

	rand.Seed(int64(seed))
	for i := range matrix {
		for j := range matrix[i] {
			randomNumber := rand.Intn(2)
			matrix[i][j] = randomNumber
		}
	}
	ClientePerf()
}

func Cliente() {
	// 1: Conectar ao servidor RPC - host/porta
	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor", err)
		return
	}
	defer client.Close()

	// 2: Invocar a operação remota
	req := impl.Request{P1: 10, P2: 20}
	//a funçao reply diz o que vai ser retornado da função add da calculadora
	rep := impl.Reply{}
	err = client.Call("Calculadora.Add", req, &rep)
	if err != nil {
		fmt.Println("Erro na chamada remota:", err)
		return
	}

	// 3: Imprimir o resultado
	fmt.Printf("Add(%v,%v) = %v \n", req.P1, req.P2, rep.R)
}

func ClientePerf() {
	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor", err)
		return
	}
	defer client.Close()

	req := ConwayGame.Request{P1: 10, P2: 20}
	rep := impl.Reply{}
	for i := 0; i < shared.StatisticSample; i++ {
		t1 := time.Now()
		for j := 0; j < shared.SampleSize; j++ {
			err = client.Call("Calculadora.Add", req, &rep)
			shared.ChecaErro(err, "Erro na invocação da Calculadora remota...")
		}
		fmt.Printf("tcp;%v\n", time.Now().Sub(t1).Milliseconds())
	}
}
