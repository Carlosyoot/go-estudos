package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func Initialize() {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.BestSpeed))

	InitializeRoutes(router)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      65 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Println("Servidor rodando em :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("erro ao subir servidor: %v", err)
	}
}
