package database

import (
	"time"

	"github.com/Carlosyoot/go-estudos/database/schema"
)

func InsertToLog(file, info, owner string) error {

	req := schema.ClassType{
		NomeRemessa: file,
		Mensagem:    info,
		Prestador:   owner,
		DataEnvio:   time.Now(),
	}

	if err := DB.Create(&req).Error; err != nil {
		return err
	}

	return nil
}
