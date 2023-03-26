package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MoveFilesForm struct {
	FromFileIDList []uint `json:"fromFileIDList"`
	ToFolderPath   string `json:"toFolderPath"`
}

func moveFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := MoveFilesForm{}
	c.BindJSON(&form)
	err := models.MoveFiles(userID.(uint), form.FromFileIDList, form.ToFolderPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
