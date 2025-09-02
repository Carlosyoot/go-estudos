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

	if err := infra.InitObserver(contexto); err != nil {
		log.Fatalf("Erro ao iniciar observador: %v", err)
	}

	server := router.InitializeRouter()

	go func() {
		log.Println("Servidor rodando em", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	<-contexto.Done()

	shutdownContexto, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownContexto); err != nil {
		log.Printf("Erro no shutdown: %v", err)
	}

	log.Println("Sinal recebido. Encerrando...")
	log.Println("Finalizado.")
}
