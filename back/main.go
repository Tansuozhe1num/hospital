package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"hospital-system/load"
	"hospital-system/server/httpserver"
	"log"
)

func main() {
	ctx := context.Background()
	load.Load(ctx)

	router := gin.Default()

	router.Use(corsMiddleware())

	httpserver.SetupRoutes(router)
	httpserver.SetUpFronted(router)
	
	log.Println("医院挂号系统后端服务启动在 http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
