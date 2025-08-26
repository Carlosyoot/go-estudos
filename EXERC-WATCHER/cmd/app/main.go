package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Carlosyoot/go-estudos/internal/infra"
	"github.com/Carlosyoot/go-estudos/router"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Variáveis de ambiente não carregadas: %v", err)
	}

	contexto, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := infra.Observar(contexto, os.Getenv("BASEPATH")); err != nil {
		log.Fatalf("Falha ao observar diretório %q: %v", os.Getenv("BASEPATH"), err)
	}

	router.Initialize()
	<-contexto.Done()

}
