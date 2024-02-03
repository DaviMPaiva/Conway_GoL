package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server on localhost:8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send a message to the server
	message := "Hello, server!"
	conn.Write([]byte(message))

	// Receive the echoed message from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// Print the echoed message
	fmt.Printf("Server response: %s\n", buffer[:n])
}
