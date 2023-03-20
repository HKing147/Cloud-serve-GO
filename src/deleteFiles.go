package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type deleteFilesForm struct {
	UserFileIDList []uint `json:"userFileIDList"`
	Path           string `json:"path"`
}

func deleteFiles(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := deleteFilesForm{}
	c.BindJSON(&form)
	err := models.DeleteFiles(userID.(uint), form.UserFileIDList, form.Path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
