package main

import (
	"log"
	"main/config"
	"main/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()
	models.Init()

	r := gin.Default()

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":8080"
	}

	initRouter(r)

	log.Printf("server is running at %s", addr)
	r.Run(addr)
}
