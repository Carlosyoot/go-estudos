package config

import (
	"fmt"
	"log"
	"time"

	oracle "github.com/dzwvip/gorm-oracle"
	goora "github.com/sijms/go-ora/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := goora.BuildUrl("localhost", 1521, "PETROSHOW", "VIASOFTGP", "VIASOFTGP", nil)

	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{

		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("Erro ao conectar no Oracle: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erro ao obter o pool de conexões: %v", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	fmt.Println("Conectado no Oracle com pool configurado!")
	DB = db
}
