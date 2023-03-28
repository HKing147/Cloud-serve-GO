package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 获取某个文件夹下的文件列表
func getFileListByFolderID(c *gin.Context) {
	folderIDStr := c.Query("folderID")
	folderID, _ := strconv.Atoi(folderIDStr)
	sortMethod := c.Query("sortMethod")
	fileList, err := models.GetFileListByFolderID(uint(folderID), sortMethod)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "fileList": fileList})
}
