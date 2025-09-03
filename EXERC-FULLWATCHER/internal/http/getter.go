package http

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Carlosyoot/go-estudos/internal/model"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html/charset"
)

type soapEnvelope struct {
	Body soapBody `xml:"Body"`
}

type soapBody struct {
	ListResp listarArquivosResp `xml:"webservice.ListarArquivosResponse"`
}

type listarArquivosResp struct {
	Return string `xml:"return"`
}

func GetIndexSoap() ([]string, error) {
	payloadStruct := model.BasicFinnetSoap{
		Servico:     "urn:ListarArquivosRequest",
		Usuario:     os.Getenv("/"),
		Senha:       os.Getenv("/"),
		CaixaPostal: os.Getenv("/"),
		Encode:      "UTF-8",
	}
	payload := model.MontarSoapListagem(payloadStruct)

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(500 * time.Millisecond)

	resp, err := client.R().
		SetHeader("Content-Type", "text/xml; charset=utf-8").
		SetHeader("SOAPAction", "urn:fml.webservice-edi.finnet.com.br#ListarArquivos").
		SetBody(payload).
		Post("https://webservice-edi.finnet.com.br")
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar soap: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("http status fora do padrão: %d", resp.StatusCode())
	}

	dec := xml.NewDecoder(bytes.NewReader(resp.Body()))
	dec.CharsetReader = charset.NewReaderLabel

	var env soapEnvelope
	if err := dec.Decode(&env); err != nil {
		return nil, fmt.Errorf("falha ao decodificar XML SOAP: %w", err)
	}

	names := splitNonEmptyLines(env.Body.ListResp.Return)
	return names, nil
}

func splitNonEmptyLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	parts := strings.Split(s, "\n")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
