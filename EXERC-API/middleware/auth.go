package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func getToken(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"Error": "Token ausente ou malformado",
		})
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

	expToken := os.Getenv("TOKEN")
	if expToken == "" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Error": "Variável de ambiente TOKEN não definida",
		})
		return
	}

	if token == "" || token != expToken {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token inválido",
		})
		return
	}

	ctx.Next()
}

func AuthBearer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getToken(ctx)
	}
}
