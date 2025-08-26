package database

import (
	"context"
	"database/sql"

	"github.com/Carlosyoot/go-estudos/config"
)

func QuerySimples() ([]map[string]any, error) {
	var out []map[string]any
	err := config.DB.
		Raw("SELECT * FROM PUSUARIO").
		Scan(&out).Error
	return out, err
}
func QuerySimplesV2(ctx context.Context) (*sql.Rows, []string, error) {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return nil, nil, err
	}

	rows, err := sqlDB.QueryContext(ctx, "SELECT * FROM PUSUARIO")
	if err != nil {
		return nil, nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		rows.Close()
		return nil, nil, err
	}

	return rows, cols, nil
}
