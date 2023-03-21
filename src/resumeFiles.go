package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type resumeFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
}

// 恢复被删除的文件
func resumeFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := resumeFilesForm{}
	c.BindJSON(&form)
	fmt.Println("userID: ", userID.(uint), "form: ", form.UserFileIDList)
	err := models.ResumeFiles(userID.(uint), form.UserFileIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
