package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CompletelyDeleteFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
}

// 从回收站中彻底删除文件
func completelyDeleteFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := CompletelyDeleteFilesForm{}
	c.BindJSON(&form)
	err := models.CompletelyDeleteFiles(userID.(uint), form.UserFileIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
