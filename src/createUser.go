package main

import (
	"Cloud-serve/src/models"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)
	// 加密密码
	h := md5.New()
	h.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))
	err := models.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}})
}
