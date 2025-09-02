package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Carlosyoot/go-estudos/internal/infra"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Variáveis de ambiente não carregadas: %v", err)
	}

	contexto, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumer := func(path string) error {
		log.Printf("Arquivo detectado (simulação de consumo): %s", path)
		return nil
	}

	if err := infra.Observar(contexto, os.Getenv("BASEPATH"), consumer); err != nil {
		log.Fatalf("Falha ao observar diretório %q: %v", os.Getenv("BASEPATH"), err)
	}

	<-contexto.Done()

	log.Println("Sinal recebido. Encerrando...")
	log.Println("Finalizado.")
}
