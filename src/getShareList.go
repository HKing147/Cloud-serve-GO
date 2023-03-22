package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getShareList(c *gin.Context) {
	userID, _ := c.Get("userID")
	fileList, err := models.GetShareList(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Meta": models.Meta{0, "success"}, "fileList": fileList})
}
