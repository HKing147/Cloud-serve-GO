package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getShareByID(c *gin.Context) {
	shareID, _ := strconv.Atoi(c.Query("shareID"))
	shareInfo, fileList, err := models.GetShareByID(uint(shareID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "fileList": fileList, "shareInfo": shareInfo})
}
