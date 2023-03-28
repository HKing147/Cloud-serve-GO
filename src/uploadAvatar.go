package main

import (
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func uploadAvatar(c *gin.Context) {
	file, _ := c.FormFile("file")

	userID, _ := c.Get("userID")
	avatar, err := OSS.UploadFile(file, "avatar/"+strconv.Itoa(int(userID.(uint)))+"_"+file.Filename)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	// 更新User表的avatar
	err = models.UpdateAvatar(userID.(uint), avatar)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "avatar": avatar})
}
