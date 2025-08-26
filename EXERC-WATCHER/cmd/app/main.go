package main

import (
	"log"

	"github.com/Carlosyoot/go-estudos/router"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Variáveis de ambiente não carregadas")
	}

	router.Initialize()

}
