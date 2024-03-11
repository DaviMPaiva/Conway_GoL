package main

import (
	"encoding/json"
	"fmt"
	"conway/rabbitmq/impl"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// cria conexão com o broker
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Não foi possível se conectar ao broker", err)
		return
	}
	defer conn.Close()

	// cria um canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Não foi possível estabelecer um canal de comunicação com o broker", err)
		return
	}
	defer ch.Close()

	// declara a fila
	q, err := ch.QueueDeclare(
		"request_queue",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("Não foi possível criar a fila no broker", err)
		return
	}

	// prepara o recebimento de mensagens do cliente
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("Falha ao registrar o consumidor no broker", err)
		return
	}

	fmt.Println("Conway pronto...")
	for d := range msgs {
		// recebe request
		msg := impl.Request{}
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			fmt.Println("Falha ao desserializar a mensagem", err)
			return
		}

		// processa request
		result := impl.ConwayGame{}.Initialize(msg)

		// prepara resposta
		replyMsg := impl.Reply{Matrix_result: result}
		replyMsgBytes, err := json.Marshal(replyMsg)
		if err != nil {
			fmt.Println( "Falha ao serializar mensagem", err)
			return
		}

		// publica resposta
		err = ch.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId, // usa correlation id do request
				Body:          replyMsgBytes,
			},
		)
		if err != nil {
			fmt.Println( "Falha ao enviar a mensagem para o broker", err)
			return
		}
	}
}
