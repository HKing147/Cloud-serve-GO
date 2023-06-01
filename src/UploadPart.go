package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type UploadPartForm struct {
	File     *multipart.FileHeader `form:"file"`
	Idx      int                   `form:"idx"`
	UploadID string                `form:"UploadID"`
}

// 上传分片
func UploadPart(c *gin.Context) {
	para := UploadPartForm{}
	err := c.ShouldBind(&para)
	if err != nil {
		fmt.Println("UploadPart err:", err)
		os.Exit(-1)
	}

	filePtr, err := para.File.Open()
	if err != nil {
		fmt.Println("UploadPart err:", err)
		os.Exit(-1)
	}
	defer filePtr.Close()
	//rd := DB.ConnRedis()
	imur := InitiateMultipartUploadResult{}
	//rd.Get(para.UploadID).Scan(&imur)
	err = DB.Get(para.UploadID).Scan(&imur)
	if err != nil {
		fmt.Println("UploadPart err:", err)
		os.Exit(-1)
	}
	bucket := OSS.BucketConn()
	part, err := bucket.UploadPart(oss.InitiateMultipartUploadResult(imur), filePtr, para.File.Size, para.Idx)
	if err != nil {
		fmt.Println("UploadPart err:", err)
		os.Exit(-1)
	}
	//exist, _ := rd.Exists(para.UploadID + "_parts").Result() // 查询key是否之前存在，不存在就为这个key设置过期时间
	exist, _ := DB.Exists(para.UploadID + "_parts").Result() // 查询key是否之前存在，不存在就为这个key设置过期时间

	part1 := uploadPart(part)
	//rd.SAdd(para.UploadID+"_parts", &part1)
	DB.SAdd(para.UploadID+"_parts", &part1)
	if exist == 0 { // 第一次加入，设置过期时间
		//rd.Expire(para.UploadID+"_parts", 100*time.Second) // 保留100s
		DB.Expire(para.UploadID+"_parts", 24*time.Hour) // 保留1天
	}

	c.JSON(http.StatusOK, gin.H{
		"meta": models.Meta{0, "UploadPart success"},
	})
}
