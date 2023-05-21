package main

import (
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getUserListForm struct {
	CurrentPage int `form:"currentPage"`
	PageSize    int `form:"pageSize"`
}

func getUserList(c *gin.Context) {
	form := getUserListForm{}
	c.BindQuery(&form)
	fmt.Println(form)
	userList, totalPage, err := models.SelectAllUser(form.CurrentPage, form.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "userList": userList, "totalPage": totalPage})
}
