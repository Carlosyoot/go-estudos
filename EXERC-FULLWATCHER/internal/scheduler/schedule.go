package scheduler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	soaphttp "github.com/Carlosyoot/go-estudos/internal/http"
)

const (
	interval   = 15 * time.Minute
	runTimeout = 10 * time.Minute
	logDir     = "logs/soap"
)

func SchedulingJob(contexto context.Context) {
	t := time.NewTicker(interval)
	defer t.Stop()

	sem := make(chan struct{}, 1)

	for {
		select {
		case <-contexto.Done():
			return
		case <-t.C:
			select {
			case sem <- struct{}{}:
				go runOnce(contexto, sem)
			default:
			}
		}
	}
}

func runOnce(contexto context.Context, sem chan struct{}) {
	defer func() { <-sem }()

	ctx, cancel := context.WithTimeout(contexto, runTimeout)
	defer cancel()

	names, err := soaphttp.GetIndexSoap()
	if err != nil {
		_ = logSoapError("ListarArquivos", err.Error(), "")
		return
	}

	_ = ctx

	for _, n := range names {
		fmt.Println("arquivo:", n)
	}
}

func logSoapError(operacao, motivo, respXML string) error {
	if err := ensureDir(logDir); err != nil {
		return fmt.Errorf("erro ao garantir pasta de log: %w", err)
	}

	filename := time.Now().Format("2006-01-02_15-04-05") + "_" + operacao + "_err.log"
	path := filepath.Join(logDir, filename)

	msg := strings.TrimSpace(motivo)
	if msg == "" {
		msg = "erro não informado"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[%s]\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&b, "operacao: %s\n", operacao)
	fmt.Fprintf(&b, "motivo: %s\n\n", msg)

	if respXML != "" {
		b.WriteString("=== SOAP RESPONSE ===\n")
		b.WriteString(respXML)
		b.WriteByte('\n')
	}

	out := b.String()
	if len(out) > 1<<20 {
		out = out[:1<<20] + "\n...(truncado)\n"
	}

	return os.WriteFile(path, []byte(out), 0644)
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
