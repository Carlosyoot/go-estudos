package router

import (
	"github.com/Carlosyoot/go-estudos/handler"
	"github.com/Carlosyoot/go-estudos/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api")
	{
		v1.GET("/logs", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"logs": "nenhum log"})
		})

		v1.GET("/usuarios", middleware.AuthBearer(), handler.GetHandler)
		v1.GET("/usuariosV2", middleware.AuthBearer(), handler.GetHandlerStreamChunk)

	}
}
