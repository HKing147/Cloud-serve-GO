package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type completeUploadFileForm struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	MD5      string `json:"MD5"`
}

func completeUploadFile(c *gin.Context) {
	form := completeUploadFileForm{}
	c.BindJSON(&form)
	fmt.Println(form)
	objectName := form.Type + "/" + form.MD5 + "." + form.Type
	fileUrl := "http://" + OSS.BucketName + "." + OSS.Endpoint + "/" + objectName
	file, err := models.InsertFile(fileUrl, form.Size, form.Type, form.MD5)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	// 将文件MD5与文件ID写入Redis（永久有效）
	DB.Set(form.MD5, file.ID, 0)
	userID, _ := c.Get("userID")
	_, err = models.InsertUserFile(userID.(uint), file.ID, form.FileName, form.Path, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
