package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取相册
func getAlbum(c *gin.Context) {
	userID, _ := c.Get("userID")
	imgList, err := models.GetAlbum(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "imgList": imgList})
}
