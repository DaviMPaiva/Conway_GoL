package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	// Connect to the server on localhost:8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	//Começa a marcar o tempo
	startTime := time.Now()

	dim := 50
	board_size := 10
	epochs := 10
	// Send a message to the server
	message := strconv.Itoa(dim) + "," + strconv.Itoa(board_size) + "," + strconv.Itoa(epochs)

	conn.Write([]byte(message))

	// Receive the echoed message from the server
	buffer := make([]byte, dim*dim)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	//Calcul o tempo decorrido
	elapsedTime := time.Since(startTime)

	result := buffer[:n]
	for i := 0; i < dim; i++ {
		fmt.Printf("%s\n", result[i*dim:(i+1)*dim])
	}
	// Print the echoed message
	fmt.Printf("Tempo decorrido: %s\n", elapsedTime)
}
