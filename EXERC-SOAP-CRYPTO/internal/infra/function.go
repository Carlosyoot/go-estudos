package infra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func processarArquivosExistentes(dir string, consumer func(string) error, arquivoDir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.Contains(path, "processados") {
			return nil
		}

		if strings.EqualFold(filepath.Ext(path), ".REM") {

			go processarArquivo(path, consumer, arquivoDir)
		}

		return nil
	})
}

func processarArquivo(filePath string, consumer func(string) error, arquivoDir string) {
	time.Sleep(50 * time.Millisecond)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	// ============= CONTROLE DE DUPLICATAS =============
	leitor.Lock()
	if _, exists := MapValue[filePath]; exists {
		leitor.Unlock()
		fmt.Printf("Arquivo %s já está sendo processado\n", filepath.Base(filePath))
		return
	}

	MapValue[filePath] = struct{}{}
	leitor.Unlock()

	// ============= CONSUMO =============
	fmt.Printf("Consumindo arquivo: %s\n", filepath.Base(filePath))

	if err := consumer(filePath); err != nil {
		fmt.Printf("Erro ao consumir %s: %v\n", filepath.Base(filePath), err)

		leitor.Lock()
		delete(MapValue, filePath)
		leitor.Unlock()
		return
	}

	nomeArquivo := filepath.Base(filePath)
	destino := filepath.Join(arquivoDir, nomeArquivo)

	if err := clonarArquivo(filePath, destino); err != nil {
		fmt.Printf("Erro ao clonar %s: %v\n", nomeArquivo, err)
	} else {
		fmt.Printf("Arquivo clonado: %s\n", nomeArquivo)
	}

	// ============= LIMPEZA =============
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Erro ao remover arquivo original %s: %v\n", nomeArquivo, err)
	}

	leitor.Lock()
	delete(MapValue, filePath)
	leitor.Unlock()
}

func clonarArquivo(origem, destino string) error {
	src, err := os.Open(origem)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
