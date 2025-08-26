package handler

import (
	"net/http"

	"github.com/Carlosyoot/go-estudos/internal/infra"
	"github.com/gin-gonic/gin"
)

func WatchHandler(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"Total": infra.Contagem(),
	})
}
