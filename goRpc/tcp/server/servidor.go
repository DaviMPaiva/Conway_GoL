package main

import (
	"conway/goRpc/impl"
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	// 1: Criar instância da calculadora.
	conwaygameService := new(impl.ConwayGame)

	// 2: Registrar a instância da calculadora no RPC
	server := rpc.NewServer()
	err := server.Register(conwaygameService)
	if err != nil {
		fmt.Println("Erro ao conectar o Conway Game", err)
		return
	}

	// 3: Criar listener para as conexões remotas
	listener, err := net.Listen("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao iniciar o Conway Game", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Servidor RPC pronto (RPC-TCP) na porta %v...\n", 1313)

	// 4: Aceitar e processar requisições remotas
	server.Accept(listener)
}
