package handler

import (
	"github.com/Carlosyoot/go-estudos/internal/infra"
	"github.com/gin-gonic/gin"
)

func ContagemHandler(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"Rems processados": infra.Contagem(),
	})
}
