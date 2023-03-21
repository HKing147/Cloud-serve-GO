package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func searchFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	fileName := c.Query("fileName")
	fileList, err := models.SearchFile(userID.(uint), fileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "fileList": fileList})
}
