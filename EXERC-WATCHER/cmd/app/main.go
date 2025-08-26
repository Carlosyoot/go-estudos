package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := router.NewServer()
	go func() {
		log.Println("Servidor rodando em", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	<-contexto.Done()

	log.Println("Sinal recebido. Encerrando...")

	shutdownContexto, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownContexto); err != nil {
		log.Printf("Erro no shutdown: %v", err)
	}

	log.Println("Finalizado.")
}
