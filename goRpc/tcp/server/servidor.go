package main

import (
	"conway/goRpc/impl"
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	conwaygameService := new(impl.ConwayGame)

	server := rpc.NewServer()
	err := server.Register(conwaygameService)
	if err != nil {
		fmt.Println("Erro ao conectar o Conway Game", err)
		return
	}

	listener, err := net.Listen("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao iniciar o Conway Game", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Servidor RPC pronto (RPC-TCP) na porta %v...\n", 1313)

	server.Accept(listener)
}
