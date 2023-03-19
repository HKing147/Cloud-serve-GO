package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CollectedFilesForm struct {
	FileIDList []uint `json:"fileIDList"`
}

func collectedFiles(c *gin.Context) {
	//fileIDList, _ := c.GetPostFormArray("fileIDList")
	form := CollectedFilesForm{}
	c.BindJSON(&form)
	//fileIDList, _ := c.GetPostForm("fileIDList")
	userID, _ := c.Get("userID")
	fmt.Println("fileIDList: ", form.FileIDList)
	err := models.CollectedFile(userID.(uint), form.FileIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
