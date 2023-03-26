package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type shareFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
}

func shareFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := shareFilesForm{}
	c.BindJSON(&form)
	err := models.InsertShare(userID.(uint), form.UserFileIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
