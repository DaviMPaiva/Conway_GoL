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

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:1313")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	for i := 0; i < 10000; i++ {

		message := strconv.Itoa(dim) + "," + strconv.Itoa(board_size) + "," + strconv.Itoa(epochs) + "," + strconv.Itoa(seed)

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to UDP connection:", err)
			return
		}

		buffer := make([]byte, dim*dim)
		n, _, err = conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			return
		}
		elapsedTime := time.Since(start_time)
		times = append(times, int(elapsedTime))
		fmt.Fprintf(file, "%s\n", elapsedTime)
		fmt.Println("pacote recebido")
	}
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	end_time := time.Since(start_time)

	fmt.Printf("Tempo decorrido: %s\n", end_time)
	print(n)
}
