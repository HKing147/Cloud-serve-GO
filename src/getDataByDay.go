package main

import (
	"Cloud-serve/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getDataByDay(c *gin.Context) {
	// 获取七天的数据
	t := time.Now().AddDate(0, 0, -6)
	dates := []time.Time{}
	registerCntList, shareCntList, uploadCntList := []int64{}, []int64{}, []int64{}
	for i := 0; i < 7; i++ {
		y, m, d := t.Date()
		date := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
		registerCnt, err := models.GetRegisterCntByDay(date)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
			return
		}

		shareCnt, err := models.GetShareCntByDay(date)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
			return
		}

		uploadCnt, err := models.GetUploadCntByDay(date)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "error"}})
			return
		}

		dates = append(dates, date)
		registerCntList = append(registerCntList, registerCnt)
		shareCntList = append(shareCntList, shareCnt)
		uploadCntList = append(uploadCntList, uploadCnt)

		t = t.AddDate(0, 0, 1)
	}
	c.JSON(http.StatusOK, gin.H{"meta": models.Meta{0, "success"}, "dates": dates,
		"registerCntList": registerCntList, "shareCntList": shareCntList, "uploadCntList": uploadCntList})
}
