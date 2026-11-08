package main

import (
	"log"
)

func main() {
	router := SetupRouter()

	log.Println("User Service iniciado na porta 8081")

	if err := router.Run(":8081"); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
