package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getFiles(c *gin.Context) {
	userFileIDList := c.QueryArray("userFileIDList[]")
	fmt.Println(userFileIDList)
	userID, _ := c.Get("userID")
	err, fileList := models.GetFilesByUserFileIDList(userFileIDList, userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "fileList": fileList})
}
