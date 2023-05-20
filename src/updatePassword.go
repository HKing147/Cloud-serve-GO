package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type updatePasswordForm struct {
	UserID  uint   `json:"id"`
	OldPass string `json:"oldPass"`
	NewPass string `json:"newPass"`
}

func updatePassword(c *gin.Context) {
	form := updatePasswordForm{}
	c.BindJSON(&form)
	err := models.UpdatePassword(form.UserID, form.OldPass, form.NewPass)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
