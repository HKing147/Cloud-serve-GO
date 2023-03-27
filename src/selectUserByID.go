package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func selectUserByID(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("userID"))
	fmt.Println(userID)
	user, err := models.SelectUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "user": user})
}
