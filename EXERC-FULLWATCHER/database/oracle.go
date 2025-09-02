package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	oracle "github.com/dzwvip/gorm-oracle"
	goora "github.com/sijms/go-ora/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	maxOpen  = 30
	maxIdle  = 15
	lifeTime = 30 * time.Minute
	idleTime = 5 * time.Minute
	pingTO   = 2 * time.Second
)

func ConnectDatabase() error {
	dsn := goora.BuildUrl("localhost", 1521, "", "", "", nil)

	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return fmt.Errorf("abrindo GORM: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("obtendo *sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(lifeTime)
	sqlDB.SetConnMaxIdleTime(idleTime)

	ctx, cancel := context.WithTimeout(context.Background(), pingTO)
	defer cancel()
	if err := ping(ctx, sqlDB); err != nil {
		return fmt.Errorf("ping Oracle: %w", err)
	}

	DB = db
	return nil
}

func ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}
