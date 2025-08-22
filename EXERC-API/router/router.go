package router

import "github.com/gin-gonic/gin"

func externalHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "obtido com sucesso",
	})
}

func Initialize() {

	router := gin.Default()
	router.GET("/ObterSoap", externalHandler)
	router.Run(":8080")

}
