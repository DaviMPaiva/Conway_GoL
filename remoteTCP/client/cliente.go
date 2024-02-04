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

	r, _ := net.ResolveTCPAddr("tcp", "localhost:8080")

	for i := 0; i < 10000; i++ {

		//Começa a marcar o tempo
		startTime := time.Now()

		conn, err := net.DialTCP("tcp", nil, r)

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

		err = conn.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		//Calcul o tempo decorrido
		elapsedTime := time.Since(startTime).Microseconds()

		// Imprimi a matriz

		//result := buffer[:n]
		//for i := 0; i < dim; i++ {
		//	fmt.Printf("%s\n", result[i*dim:(i+1)*dim])
		//}
		// Print the echoed message
		//fmt.Printf("Tempo decorrido: %s\n", elapsedTime)

		times = append(times, int(elapsedTime))
		fmt.Fprintf(file, "%d\n", elapsedTime)
		fmt.Printf("pacote recebido numero %d\n", i)
	}
	print(n)
}
