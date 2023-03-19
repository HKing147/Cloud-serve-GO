package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getRecycle(c *gin.Context) {
	userID, _ := c.Get("userID")
	recycle, err := models.GetRecycle(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "recycle": recycle})
}
