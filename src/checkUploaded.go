package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func checkUploaded(c *gin.Context) {
	MD5 := c.Query("MD5")
	fileIDStr, _ := DB.Get(MD5).Result()
	fmt.Println(fileIDStr)
	if fileIDStr == "" { // 没上传过
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "没上传过。"}})
		return
	}

	userID, _ := c.Get("userID")
	fileID, _ := strconv.Atoi(fileIDStr)
	fileName := c.Query("fileName")
	filePath := c.Query("path")
	models.InsertUserFile(userID.(uint), uint(fileID), fileName, filePath, false)
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "上传过。"}})
}
