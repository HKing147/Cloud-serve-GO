package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type shareFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
	ShareMethod    bool   `json:"shareMethod"`   // 0: 无密码 1: 有密码
	ShareDuration  int    `json:"shareDuration"` // 有效期（天）
}

func shareFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := shareFilesForm{}
	c.BindJSON(&form)
	share, err := models.InsertShare(userID.(uint), form.UserFileIDList, form.ShareMethod, form.ShareDuration)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	shareUrl := fmt.Sprintf("share_%v_%v", userID, share.ID)
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "shareUrl": shareUrl, "password": share.Password})
}
