package infra

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//=====================================//

var leitor sync.RWMutex
var MapValue = make(map[string]struct{})

//=====================================//

func InitObserver(contexto context.Context) error {
	dir := os.Getenv("BASEPATH")
	if dir == "" {
		return fmt.Errorf("erro ao obter o diretório base: variável de ambiente BASEPATH não definida")
	}

	processedDir, err := checkDir(dir)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório base: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(e.Name()), ".REM") {
			continue
		}

		full := filepath.Join(dir, e.Name())

		if _, statErr := os.Stat(full); statErr != nil {
			continue
		}

		ok, err := ProcessarArquivos(full, processedDir)
		if err != nil {
			fmt.Printf("erro ao processar %q: %v\n", full, err)
			continue
		}
		if !ok {
			log.Printf("arquivo %s movido para fallback (falha no envio)", filepath.Base(full))
			continue
		}
		log.Printf("ok: %s", filepath.Base(full))
	}

	return observar(contexto)
}

func Contabilizar(file, mode string) (string, error) {
	leitor.Lock()
	defer leitor.Unlock()

	switch mode {
	case "Add":
		MapValue[file] = struct{}{}
		return "Contabilizado com sucesso", nil
	case "Remove":
		delete(MapValue, file)
		return "Arquivo removido com sucesso", nil
	default:
		return "", fmt.Errorf("unknown mode: %s", mode)
	}
}

func Contagem() int {
	leitor.RLock()
	defer leitor.RUnlock()
	return len(MapValue)
}

//------------------------PRIVATE FUNCTIONS

func observar(contexto context.Context) error {

	dir := os.Getenv("BASEPATH")
	if dir == "" {
		return fmt.Errorf("erro ao obter o diretório base: variável de ambiente BASEPATH não definida")
	}

	processedDir, err := checkDir(dir)
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("erro ao iniciar observador: %v", err)
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return fmt.Errorf("erro ao adicionar diretório no observador: %v", err)
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

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {

					if err := validateTypeFile(event.Name); err != nil {
						continue
					}

					if _, statErr := os.Stat(event.Name); statErr != nil {
						continue
					}

					if _, err := Contabilizar(event.Name, "Add"); err != nil {
						log.Printf("erro ao contabilizar: %v", err)
					}

					ok, err := ProcessarArquivos(event.Name, processedDir)
					if err != nil {
						fmt.Printf("erro ao processar %q: %v\n", event.Name, err)
						continue
					}
					if !ok {
						log.Printf("arquivo %s movido para fallback (falha no envio)", filepath.Base(event.Name))
						continue
					}
					// sucesso
					log.Printf("ok: %s", filepath.Base(event.Name))

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher errors channel fechado")
					return
				}
				if err != nil {
					fmt.Printf("erro no watcher: %v\n", err)
				}

			case <-contexto.Done():
				fmt.Println("Context cancelado, parando watcher")
				return
			}
		}
	}()

	return nil
}

func validateTypeFile(path string) error {
	if !strings.EqualFold(filepath.Ext(path), ".REM") {
		return fmt.Errorf("arquivo não é .REM: %s", path)
	}
	return nil
}

func checkDir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	ArchiveDir := filepath.Join(dir, "Processed")
	if err := os.MkdirAll(ArchiveDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de arquivo: %v", err)
	}

	return ArchiveDir, nil
}
