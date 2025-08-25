package handler

import (
	"net/http"

	"github.com/Carlosyoot/go-estudos/database"
	"github.com/gin-gonic/gin"
)

func GetHandler(ctx *gin.Context) {
	data, err := database.QuerySimples()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "falha ao buscar dados",
			"detail": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
