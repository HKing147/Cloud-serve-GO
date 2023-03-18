package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/middleware"
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
)

func main() {
	models.Init()
	DB.InitRedis()
	models.InsertFile("", 0, "folder", "") // 文件夹 fileID为1
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.AuthMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/upload", uploadFile)
	r.GET("/api/InitiateMultipartUpload", InitiateMultipartUpload)
	r.POST("/api/UploadPart", UploadPart)
	r.POST("/api/CompleteMultipartUpload", CompleteMultipartUpload)
	r.POST("/api/login", models.Login)
	r.POST("/api/register", models.Register)
	r.GET("/api/getCheckCode", getCheckCode)
	r.GET("/api/getFileList", models.GetFileList)
	r.GET("/api/checkUploaded", checkUploaded)
	r.POST("/api/createFolder", createFolder)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
