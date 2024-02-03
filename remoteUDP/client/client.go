package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start measuring time
	startTime := time.Now()

	dim := 50
	board_size := 10
	epochs := 10
	seed := 42

	// Send a message to the server
	message := strconv.Itoa(dim) + "," + strconv.Itoa(board_size) + "," + strconv.Itoa(epochs) + "," + strconv.Itoa(seed)

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to UDP connection:", err)
		return
	}

	// Receive the echoed message from the server
	buffer := make([]byte, dim*dim)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading from UDP connection:", err)
		return
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	result := buffer[:n]
	for i := 0; i < dim; i++ {
		fmt.Printf("%s\n", result[i*dim:(i+1)*dim])
	}

	// Print elapsed time
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
