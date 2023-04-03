package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteSharesForm struct {
	ShareIDList []uint `json:"shareIDList"`
}

func deleteShares(c *gin.Context) {
	userID, _ := c.Get("userID")
	form := DeleteSharesForm{}
	err := c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	err = models.DeleteShares(userID.(uint), form.ShareIDList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
