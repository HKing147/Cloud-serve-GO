package main

import (
	"Cloud-serve/src/CORS"
	"Cloud-serve/src/DB"
	"github.com/gin-gonic/gin"
)

func main() {
	DB.ConnRedis()
	r := gin.Default()
	r.Use(CORS.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/upload", uploadFile)
	r.GET("/api/InitiateMultipartUpload", InitiateMultipartUpload)
	r.POST("/api/UploadPart", UploadPart)
	r.POST("/api/CompleteMultipartUpload", CompleteMultipartUpload)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
