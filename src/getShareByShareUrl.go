package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getShareByShareUrl(c *gin.Context) {
	shareUrl := c.Query("shareUrl")
	sortMethod := c.Query("sortMethod")
	shareInfo, fileList, err := models.GetShareByShareUrl(shareUrl, sortMethod)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "shareInfo": shareInfo, "fileList": fileList})
}
