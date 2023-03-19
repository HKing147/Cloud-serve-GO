package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getCollectedList(c *gin.Context) {
	userID, _ := c.Get("userID")
	collectedList, err := models.GetCollectedList(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "collectedList": collectedList})
}
