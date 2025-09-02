package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// ============= VARIÁVEIS GLOBAIS =============

// Vários leitores ou um escritor
var leitor = sync.RWMutex{}

// struct struct não consome bytes, apenas seta valor
var MapValue = make(map[string]struct{})

// ============= FUNÇÕES =============

func Contagem() int {
	leitor.RLock()
	defer leitor.RUnlock()
	return len(MapValue)
}

func Observar(ctx context.Context, dir string, consumer func(string) error) error {

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	arquivoDir := filepath.Join(dir, "processados")
	if err := os.MkdirAll(arquivoDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de arquivo: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("erro ao criar watcher: %v", err)
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return fmt.Errorf("erro ao adicionar diretório ao watcher: %v", err)
	}

	if err := processarArquivosExistentes(dir, consumer, arquivoDir); err != nil {
		watcher.Close()
		return fmt.Errorf("erro ao processar arquivos existentes: %v", err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {

			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("Watcher events channel fechado")
					return
				}

				if event.Has(fsnotify.Create) {
					if strings.EqualFold(filepath.Ext(event.Name), ".REM") {
						go processarArquivo(event.Name, consumer, arquivoDir)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher errors channel fechado")
					return
				}
				if err != nil {
					fmt.Printf("Erro no watcher: %v\n", err)
				}

			case <-ctx.Done():
				fmt.Println("Context cancelado, parando watcher")
				return
			}
		}
	}()

	return nil
}
