package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type renameFileForm struct {
	UserFileID  uint   `json:"userFileID"`
	NewFileName string `json:"newFileName"`
}

func renameFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := renameFileForm{}
	c.BindJSON(&form)
	err := models.RenameFile(userID.(uint), form.UserFileID, form.NewFileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
