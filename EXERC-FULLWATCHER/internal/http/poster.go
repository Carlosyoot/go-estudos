package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func SenderSoap(xml string) (string, error) {
	payload := strings.TrimSpace(xml)
	if payload == "" {
		return "", fmt.Errorf("xml soap está vazio")
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(500 * time.Millisecond)

	resp, err := client.R().
		SetHeader("Content-Type", "text/xml; charset=utf-8").
		SetHeader("SOAPAction", "urn:fml.webservice-edi.finnet.com.br#EnviarArquivos").
		SetBody(payload).
		Post("https://webservice-edi.finnet.com.br")
	if err != nil {
		return "", fmt.Errorf("erro ao enviar soap: %v", err)
	}

	body := string(resp.Body())

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return "", fmt.Errorf("http status fora do padrão: %d", resp.StatusCode())
	}

	if !strings.Contains(strings.ToLower(body), "executado com sucesso") {
		trimmed := body
		if len(trimmed) > 1024 {
			trimmed = trimmed[:1024] + "...(truncado)"
		}
		return "", fmt.Errorf("envio não confirmado: %s", trimmed)
	}

	return body, nil
}
