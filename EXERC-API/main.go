package main

import (
	"log"

	"github.com/Carlosyoot/go-estudos/config"
	"github.com/Carlosyoot/go-estudos/router"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: .env não encontrado (ok em produção se variáveis já estiverem no ambiente)")
	}

	config.ConnectDatabase()
	router.Initialize()

}
