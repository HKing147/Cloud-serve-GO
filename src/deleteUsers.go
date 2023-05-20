package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func deleteUsers(c *gin.Context) {
	userIDList := []uint{}
	c.BindJSON(&userIDList)
	fmt.Println("userIDList: ", userIDList)
	err := models.DeleteUsers(userIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
