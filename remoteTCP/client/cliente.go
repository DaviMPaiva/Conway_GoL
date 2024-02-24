package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Size   int     `json:"size"`
	Matrix [][]int `json:"matrix"`
}

func main() {
	dim, _ := strconv.Atoi(os.Args[1])
	epochs, _ := strconv.Atoi(os.Args[2])
	n := 0
	file, _ := os.OpenFile("../../outputs/TCP_"+os.Args[1]+"_"+os.Args[2]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)

	r, _ := net.ResolveTCPAddr("tcp", "localhost:8080")

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
	data := Data{
		Size:   dim,
		Matrix: matrix,
	}

	//se conecta com o servidor
	conn, err := net.DialTCP("tcp", nil, r)

	for i := 0; i < int(epochs); i++ {

		//Começa a marcar o tempo
		startTime := time.Now()

		//prepara os dados
		bytes_men, _ := json.Marshal(data)

		// manda mensagem para o servidor
		conn.Write([]byte(bytes_men))

		// recebe a mensagem do servidor
		buffer := make([]byte, 65535)
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		//Calcula o tempo decorrido
		elapsedTime := time.Since(startTime).Microseconds()
		fmt.Fprintf(file, "%d\n", elapsedTime)

		//desserializa
		var receivedData Data
		err = json.Unmarshal(buffer[:n], &receivedData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			//fmt.Println("Unmarshaled data:", receivedData)
			data = receivedData
		}

		//fmt.Printf("Tempo decorrido: %s\n", elapsedTime)
		fmt.Printf("pacote recebido numero %d\n", i)
	}
	//fecha a conexao
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	print(n)
}
