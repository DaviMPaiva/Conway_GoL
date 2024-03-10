package main

import (
	"conway/rabbitmq/impl"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	dim, _ := strconv.Atoi(os.Args[1])
	epochs, _ := strconv.Atoi(os.Args[2])
	print_result, _ := strconv.Atoi(os.Args[3])
	display, _ := strconv.Atoi(os.Args[4])
	file, _ := os.OpenFile("../../../outputs/GoRPC_"+os.Args[1]+"_"+os.Args[2]+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0222)

	matrix := make([][]int, dim)
	//prepare matrix
	if display > 0 {
		matrix = plainTextReader(dim)
	} else {
		for i := range matrix {
			matrix[i] = make([]int, dim)
		}
		//rand init
		rand.Seed(int64(42))
		for i := range matrix {
			for j := range matrix[i] {
				randomNumber := rand.Intn(2)
				matrix[i][j] = randomNumber
			}
		}
	}

	// conecta ao broker
	matrix_aux := matrix
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Não foi possível se conectar ao servidor de mensageria", err)
		return
	}
	defer conn.Close()

	// cria o canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Não foi possível estabelecer um canal de comunicação com o servidor de mensageria", err)
		return
	}

	defer ch.Close()

	// declara a fila para as respostas
	replyQueue, _ := ch.QueueDeclare(
		"response_queue",
		false,
		false,
		true,
		false,
		nil,
	)

	// cria servidor da fila de response
	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("Falha ao registrar o servidor no broker", err)
		return
	}

	for i := 0; i < int(epochs); i++ {
		// prepara mensagem
		msgRequest := impl.Request{Matrix: matrix_aux, Dim: dim}
		msgRequestBytes, err := json.Marshal(msgRequest)
		if err != nil {
			fmt.Println("Falha ao serializar a mensagem", err)
			return
		}

		correlationID := RandomString(32)
		start_time := time.Now()
		err = ch.Publish(
			"",
			"request_queue",
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: correlationID,
				ReplyTo:       replyQueue.Name,
				Body:          msgRequestBytes,
			},
		)

		// recebe mensagem do servidor de mensageria
		m := <-msgs
		elapsedTime := time.Since(start_time).Microseconds()
		// deserializada e imprime mensagem na tela
		msgResponse := impl.Reply{}
		err = json.Unmarshal(m.Body, &msgResponse)
		if err != nil {
			fmt.Println("Erro na deserialização da resposta", err)
			return
		}
		fmt.Fprintf(file, "%d\n", elapsedTime)
		// atualiza a nova matrix
		matrix_aux = msgResponse.Matrix_result
		// espera um tempo para printar, limpa o terminal e chama a funcao para printar
		if int(print_result) > 0 {
			time.Sleep(time.Second/5)
			fmt.Println("\033[H\033[2J")
			displayBoard(matrix_aux)
			fmt.Printf("\n\npacote recebido numero %d\n\n", i)
		} else {
			fmt.Printf("pacote recebido numero %d\n", i)
		}
	}
}

const (
	Dead  = 0
	Alive = 1
)

func displayBoard(board [][]int) {
	for _, row := range board {
		for _, val := range row {
			if val == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print("██")
			}
		}
		fmt.Println()
	}
}

func plainTextReader(dim int) [][]int {
	// Input plaintext
	// Read file
	data, err := os.ReadFile("../../../Gosper-glider_gun.txt")
	if err != nil {
		fmt.Println(err)
	}

	// Convert to string
	plaintext := string(data)

	// Split into lines
	lines := strings.Split(plaintext, "\n")

	// Width is length of first line
	width := dim

	// Height is number of lines
	height := dim

	// Create matrix
	matrix := make([][]int, height)

	for i := range matrix {
		matrix[i] = make([]int, width)
	}

	// Parse each line into matrix
	for y, line := range lines {
		for x, c := range line {
			if c == 'O' {
				matrix[y][x] = 1
			}
		}
	}

	return matrix
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
