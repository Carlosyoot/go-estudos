package schema

import "time"

type ClassType struct {
	NomeRemessa string    `gorm"colum:NOME_REMESSA"`
	Mensagem    string    `gorm"colum:MENSAGEM"`
	Prestador   string    `gorm"colum:PRESTADOR"`
	DataEnvio   time.Time `gorm:"colum:DATA_ENVIO`
}

func (ClassType) TableName() string {
	return "U_FINNET_LOGS"
}
