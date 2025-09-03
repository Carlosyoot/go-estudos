package infra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Carlosyoot/go-estudos/internal/crypto"
	"github.com/Carlosyoot/go-estudos/internal/http"
	"github.com/Carlosyoot/go-estudos/internal/model"
)

func ProcessarArquivos(file, processedDir string) (bool, error) {
	fallbackDir := os.Getenv("BASEPATHFALLBACK")
	if fallbackDir == "" {
		return false, fmt.Errorf("erro ao obter o diretório de fallback: variável de ambiente BASEPATH-FALLBACK não definida")
	}

	if err := waitFileReady(file, 10*time.Second, 200*time.Millisecond); err != nil {
		return false, fmt.Errorf("arquivo indisponível: %v", err)
	}

	md5Hash, err := crypto.MD5Encoder(file)
	if err != nil {
		return false, fmt.Errorf("erro ao calcular md5: %v", err)
	}

	b64, err := crypto.Base64Encoder(file)
	if err != nil {
		return false, fmt.Errorf("erro ao calcular base64: %v", err)
	}

	reqBody := model.FinnetSoap{
		Servico:     "urn:EnviarArquivosRequest",
		Usuario:     os.Getenv("/"),
		Senha:       os.Getenv("/"),
		CaixaPostal: os.Getenv("/"),
		Hash:        md5Hash,
		Filename:    filepath.Base(file),
		Conteudo:    b64,
		Encode:      "UTF-8",
	}

	envelope := model.MontarSoap(reqBody)

	_, sendErr := http.SenderSoap(envelope)
	if sendErr != nil {
		if err := ensureDir(fallbackDir); err != nil {
			return false, fmt.Errorf("erro ao garantir pasta de falha: %v", err)
		}
		if err := clonarArquivo(file, fallbackDir); err != nil {
			return false, fmt.Errorf("erro ao mover para pasta de falha: %v", err)
		}

		errLogPath := filepath.Join(fallbackDir, filepath.Base(file)+"_err.log")

		msg := sendErr.Error()
		var respXML string
		if idx := strings.Index(msg, "<?xml"); idx != -1 {
			respXML = msg[idx:]
			msg = strings.TrimSpace(msg[:idx])
		}

		var b strings.Builder
		fmt.Fprintf(&b, "[%s]\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(&b, "arquivo: %s\n", filepath.Base(file))
		if msg == "" {
			msg = "envio não confirmado"
		}
		fmt.Fprintf(&b, "motivo: %s\n\n", msg)

		if respXML != "" {
			b.WriteString("=== SOAP RESPONSE ===\n")
			b.WriteString(prettyXML(respXML))
			b.WriteByte('\n')
		}

		if envelope != "" {
			b.WriteString("=== SOAP REQUEST (ENVELOPE) ===\n")
			b.WriteString(prettyXML(envelope))
			b.WriteByte('\n')
		}

		out := b.String()
		if len(out) > 1<<20 {
			out = out[:1<<20] + "\n...(truncado)\n"
		}

		_ = os.WriteFile(errLogPath, []byte(out), 0644)
		return false, nil
	}

	if err := clonarArquivo(file, processedDir); err != nil {
		return false, fmt.Errorf("erro ao clonar arquivo: %v", err)
	}
	return true, nil
}

func clonarArquivo(oldDir, newDir string) error {
	dstPath := filepath.Join(newDir, filepath.Base(oldDir))

	deadline := time.Now().Add(10 * time.Second)
	var (
		content *os.File
		err     error
	)
	for {
		content, err = os.Open(oldDir)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("erro ao abrir origem: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}
	defer content.Close()

	contentDuplicate, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("erro ao criar destino: %v", err)
	}
	defer contentDuplicate.Close()

	if _, err = io.Copy(contentDuplicate, content); err != nil {
		return fmt.Errorf("erro ao copiar: %v", err)
	}
	if err = contentDuplicate.Sync(); err != nil {
		return fmt.Errorf("erro ao sincronizar: %v", err)
	}

	_ = contentDuplicate.Close()
	_ = content.Close()

	removeDeadLine := time.Now().Add(10 * time.Second)
	for {
		if err := os.Remove(oldDir); err == nil {
			break
		}
		if time.Now().After(removeDeadLine) {
			return fmt.Errorf("erro ao remover origem: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	if _, err = Contabilizar(contentDuplicate.Name(), "Add"); err != nil {
		return err
	}
	return nil
}

func waitFileReady(path string, maxWait, interval time.Duration) error {
	deadline := time.Now().Add(maxWait)
	var lastSize int64 = -1
	for {
		fi, err := os.Stat(path)
		if err == nil {
			size := fi.Size()
			f, e := os.Open(path)
			if e == nil {
				_ = f.Close()
				if size == lastSize {
					return nil
				}
				lastSize = size
			}
		}
		if time.Now().After(deadline) {
			if _, e := os.Stat(path); e != nil {
				return e
			}
			return nil
		}
		time.Sleep(interval)
	}
}

func ensureDir(dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if info.IsDir() {
			return nil
		}
		return fmt.Errorf("caminho existe mas não é diretório: %s", dir)
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
		return nil
	}
	return fmt.Errorf("erro ao verificar diretório %s: %v", dir, err)
}

func prettyXML(x string) string {
	x = strings.TrimSpace(x)
	if x == "" {
		return x
	}
	x = strings.ReplaceAll(x, "><", ">\n<")

	lines := strings.Split(x, "\n")
	var b strings.Builder
	indent := 0

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "</") && indent > 0 {
			indent--
		}

		b.WriteString(strings.Repeat("  ", indent))
		b.WriteString(l)
		b.WriteByte('\n')

		if strings.HasPrefix(l, "<") &&
			!strings.HasPrefix(l, "</") &&
			!strings.HasPrefix(l, "<?") &&
			!strings.HasSuffix(l, "/>") &&
			!strings.Contains(l, "</") {
			indent++
		}
	}
	return b.String()
}
