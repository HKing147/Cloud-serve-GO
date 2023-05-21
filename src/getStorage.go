package main

import (
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getStorage(c *gin.Context) {
	storage, err := OSS.GetStorage()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "storage": storage})
}
