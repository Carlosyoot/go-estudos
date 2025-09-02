package model

import "fmt"

type FinnetSoap struct {
	Usuario     string
	Senha       string
	CaixaPostal string
	Hash        string
	Filename    string
	Conteudo    string
	Encode      string
}

func MontarSoap(request FinnetSoap) string {
	return fmt.Sprintf(`
   <soapenv:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="/">
   <soapenv:Header/>
   <soapenv:Body>
      <urn:EnviarArquivosRequest soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
         <usuario xsi:type="xsd:string">%s</usuario>
         <senha xsi:type="xsd:string">%s</senha>
         <caixa_postal xsi:type="xsd:string">%s</caixa_postal>
         <hash xsi:type="xsd:string">%s</hash>
         <filename xsi:type="xsd:string">%s</filename>
         <conteudo xsi:type="xsd:string">%s</conteudo>
         <encode xsi:type="xsd:string">%s</encode>
      </urn:EnviarArquivosRequest>
   	</soapenv:Body>
	</soapenv:Envelope>`,
		request.Usuario,
		request.Senha,
		request.CaixaPostal,
		request.Hash,
		request.Filename,
		request.Conteudo,
		request.Encode)
}
