package main

import (
	"Cloud-serve/src/Email"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getCheckCode(c *gin.Context) {
	email := c.Query("email")
	T, _ := strconv.Atoi(c.Query("T"))
	fmt.Println(email, T)

	err := Email.SendEmail(email, T)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "验证码发送失败，请稍后重试！！！"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "验证码发送成功！"}})
}
