package database

import (
	"github.com/Carlosyoot/go-estudos/config"
)

func QuerySimples() ([]map[string]any, error) {
	var out []map[string]any
	err := config.DB.
		Raw("SELECT * FROM PUSUARIO").
		Scan(&out).Error
	return out, err
}
