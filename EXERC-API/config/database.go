package config

import (
	"fmt"
	"log"
	"time"

	"github.com/cengsin/oracle"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	url := "DATABASE/DATABASE@localhost:1521/ORLC19"

	db, err := gorm.Open(oracle.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no Oracle: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erro ao obter o pool de conexões: %v", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	fmt.Println("Conectado no Oracle com pool configurado!")

	DB = db
}
