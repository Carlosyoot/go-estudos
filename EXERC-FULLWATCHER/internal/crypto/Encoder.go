package crypto

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Base64Encoder(file string) (string, error) {
	if !strings.EqualFold(filepath.Ext(file), ".rem") {
		return "", fmt.Errorf("arquivo não é formato .REM: %s", file)
	}

	header, err := validateCNAB(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler header: %v", err)
	}

	l := len(header)
	switch {
	case l >= 230 && l <= 250:
	case l >= 390 && l <= 410:
	default:
		return "", fmt.Errorf("tamanho de header inesperado: %d", l)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %v", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func MD5Encoder(file string) (string, error) {
	if !strings.EqualFold(filepath.Ext(file), ".REM") {
		return "", fmt.Errorf("arquivo nao é formato .REM: %s", file)
	}

	header, err := validateCNAB(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler header: %v", err)
	}

	l := len(header)
	switch {
	case l >= 230 && l <= 250:
	case l >= 390 && l <= 410:
	default:
		return "", fmt.Errorf("tamanho de header inesperado: %d", l)
	}

	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer f.Close()

	hashMD5 := md5.New()
	if _, err := io.Copy(hashMD5, f); err != nil {
		return "", fmt.Errorf("erro ao calcular md5: %v", err)
	}

	sum := hashMD5.Sum(nil)
	return hex.EncodeToString(sum), nil
}

func validateCNAB(path string) (string, error) {
	cnabFILE, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer cnabFILE.Close()

	reader := bufio.NewReader(cnabFILE)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")

		if line != "" {
			return line, nil
		}

		if err == io.EOF {
			if line == "" {
				return "", fmt.Errorf("arquivo vazio ou sem header em formato cnab")
			}
			return line, nil
		}
		if err != nil {
			return "", fmt.Errorf("erro ao ler linha: %v", err)
		}
	}
}
