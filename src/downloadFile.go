package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type downloadFileForm struct {
	FileID uint `json:"id"`
}

func downloadFile(c *gin.Context) {
	form := downloadFileForm{}
	c.BindJSON(&form)
	err := models.AddDownCount(form.FileID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
