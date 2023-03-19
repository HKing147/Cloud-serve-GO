package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateFolderForm struct {
	FolderName string `json:"folderName"`
	Path       string `json:"path"`
}

func createFolder(c *gin.Context) {
	//models.InsertFile("", 0, "", "")
	form := CreateFolderForm{}
	c.BindJSON(&form)
	userID, _ := c.Get("userID")
	//folderName, _ := c.GetPostForm("folderName")
	//path, _ := c.GetPostForm("path")
	// 先创建文件夹
	folder, _ := models.InsertFile("", 0, "folder", "")
	// 再将文件夹与User关联
	err := models.InsertUserFile(userID.(uint), folder.ID, form.FolderName, form.Path, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "创建文件夹失败！！！"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
