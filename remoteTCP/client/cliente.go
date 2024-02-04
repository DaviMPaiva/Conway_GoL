package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	var times []int
	dim, _ := strconv.Atoi(os.Args[1])
	board_size, _ := strconv.Atoi(os.Args[2])
	epochs, _ := strconv.Atoi(os.Args[3])
	seed, _ := strconv.Atoi(os.Args[4])
	n := 0
	file, _ := os.OpenFile("../../outputs/"+os.Args[1]+"_"+os.Args[2]+"_"+os.Args[3]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)
	start_time := time.Now()

	for i := 0; i < 10000; i++ {

		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Error connecting:", err)
			os.Exit(1)
		}
		defer conn.Close()

		//Começa a marcar o tempo
		startTime := time.Now()

		// Send a message to the server
		message := strconv.Itoa(dim) + "," + strconv.Itoa(board_size) + "," + strconv.Itoa(epochs) + "," + strconv.Itoa(seed)

		conn.Write([]byte(message))

		// Receive the echoed message from the server
		buffer := make([]byte, dim*dim)
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		//Calcul o tempo decorrido
		elapsedTime := time.Since(startTime)

		// Imprimi a matriz

		//result := buffer[:n]
		//for i := 0; i < dim; i++ {
		//	fmt.Printf("%s\n", result[i*dim:(i+1)*dim])
		//}
		// Print the echoed message
		//fmt.Printf("Tempo decorrido: %s\n", elapsedTime)

		times = append(times, int(elapsedTime))
		fmt.Fprintf(file, "%s\n", elapsedTime)
	}
	end_time := time.Since(start_time)

	fmt.Printf("Tempo decorrido: %s\n", end_time)
	print(n)

}
