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
	file, _ := os.OpenFile("../../outputs/UDP_"+os.Args[1]+"_"+os.Args[2]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)

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
	for i := 0; i < int(epochs); i++ {

		//começa contar o tempo
		start_time := time.Now()

		//prepara os dados
		bytes_men, _ := json.Marshal(data)

		//envia mensagem
		_, err = conn.Write([]byte(bytes_men))
		if err != nil {
			fmt.Println("Error writing to UDP connection:", err)
			return
		}

		//recebe a mensagem
		buffer := make([]byte, 65535)
		n, _, err = conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			return
		}

		//calcula tempo decorrido
		elapsedTime := time.Since(start_time).Microseconds()
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
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	print(n)
}
