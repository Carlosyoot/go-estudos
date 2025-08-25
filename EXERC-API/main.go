package main

import (
	"github.com/Carlosyoot/go-estudos/config"
	"github.com/Carlosyoot/go-estudos/router"
)

func main() {

	router.Initialize()
	config.ConnectDatabase()
}
