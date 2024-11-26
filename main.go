package main

import (
	"OneDisk/lib/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Initialize() // 先初始化日志
	log.AppStart()   // 打印应用启动日志

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
