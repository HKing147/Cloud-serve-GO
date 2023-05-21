package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getFileCategory(c *gin.Context) {
	fileCategoryList, err := models.GetFileCategory()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "fileCategoryList": fileCategoryList})
}
