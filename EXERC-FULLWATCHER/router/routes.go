package router

import (
	"github.com/Carlosyoot/go-estudos/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	v1 := router.Group("/api")
	{
		v1.GET("/contagem", handler.ContagemHandler)

	}
}
