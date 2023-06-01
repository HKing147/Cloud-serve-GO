package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func updateUserInfo(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)
	fmt.Println(user)
	err := models.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
