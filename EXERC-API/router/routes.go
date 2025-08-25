package router

import "github.com/gin-gonic/gin"

func InitializeRoutes(router *gin.Engine) {
	v1 := router.Group("viasoft/api/")
	{
		v1.GET("/logs", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"logs": "nenhum",
			})
		})

	}
}
