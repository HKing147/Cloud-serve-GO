package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SaveFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
	SavePath       string `json:"savePath"`
}

// 转存别人分享的文件
func saveFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := SaveFilesForm{}
	c.BindJSON(&form)
	err := models.SaveFiles(userID.(uint), form.UserFileIDList, form.SavePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
