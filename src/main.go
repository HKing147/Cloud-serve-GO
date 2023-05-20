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
	r.POST("/api/createUser", createUser)
	r.POST("/api/deleteUsers", deleteUsers)
	r.GET("/api/getUserInfo", getUserInfo)
	r.POST("/api/updatePassword", updatePassword)
	r.POST("/api/updateUserInfo", updateUserInfo)
	r.POST("/api/uploadAvatar", uploadAvatar)
	r.GET("/api/selectUserByID", selectUserByID) // 不需要登陆
	r.GET("/api/getCheckCode", getCheckCode)
	r.GET("/api/getFileList", models.GetFileList)
	r.GET("/api/getFileListByFolderID", getFileListByFolderID)
	r.GET("/api/getUserList", getUserList)
	r.GET("/api/checkUploaded", checkUploaded)
	r.POST("/api/createFolder", createFolder)
	r.POST("/api/collectedFiles", collectedFiles)
	r.GET("/api/getCollectedList", getCollectedList)
	r.GET("/api/getAlbum", getAlbum)
	r.POST("/api/deleteFiles", deleteFiles)
	r.POST("/api/downloadFile", downloadFile)
	r.GET("/api/getRecycle", getRecycle)
	r.POST("/api/clearRecycle", clearRecycle)
	r.POST("/api/resumeFiles", resumeFiles)
	r.GET("/api/searchFile", searchFile)
	r.POST("/api/renameFile", renameFile)
	r.POST("/api/shareFiles", shareFiles)
	r.POST("/api/deleteShares", deleteShares)
	r.GET("/api/getShareByID", getShareByID)
	r.GET("/api/getShareList", getShareList)
	r.POST("/api/completelyDeleteFiles", completelyDeleteFiles)
	r.POST("/api/moveFiles", moveFiles)
	r.GET("/api/getShareByShareUrl", getShareByShareUrl)
	r.POST("/api/saveFiles", saveFiles)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
