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
	var times []int
	dim, _ := strconv.Atoi(os.Args[1])
	board_size, _ := strconv.Atoi(os.Args[2])
	epochs, _ := strconv.Atoi(os.Args[3])
	seed, _ := strconv.Atoi(os.Args[4])
	n := 0
	file, _ := os.OpenFile("../../outputs/"+string(dim)+"_"+string(board_size)+"_"+string(epochs)+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)

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
	data := Data{
		Size:   dim,
		Matrix: matrix,
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:5151")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	fmt.Println("Conexão estabelecida. Iniciando envio de pacotes")
	for i := 0; i < 2; i++ {
		start_time := time.Now()

		//prepara os dados
		bytes_men, _ := json.Marshal(data)

		_, err = conn.Write([]byte(bytes_men))
		if err != nil {
			fmt.Println("Error writing to UDP connection:", err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err = conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			return
		}

		var receivedData Data
		err = json.Unmarshal(buffer[:n], &receivedData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			fmt.Println("Unmarshaled data:", receivedData)
			data = receivedData
		}

		elapsedTime := time.Since(start_time).Microseconds()
		times = append(times, int(elapsedTime))
		fmt.Fprintf(file, "%d\n", elapsedTime)
		fmt.Printf("pacote recebido numero %d\n", i)
	}
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	print(n)
}
