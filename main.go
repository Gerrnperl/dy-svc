package main

import (
	"log"
	"main/config"
	"main/models"

	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()
	models.Init()

	r := gin.Default()

	addr := config.Address + ":" + config.Port

	initRouter(r)

	log.Printf("server is running at %s", addr)
	r.Run(addr)
}
