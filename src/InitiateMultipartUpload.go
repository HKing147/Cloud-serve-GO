package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type InitiateMultipartUploadResult oss.InitiateMultipartUploadResult

func (o *InitiateMultipartUploadResult) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *InitiateMultipartUploadResult) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, o)
}

type uploadPart oss.UploadPart

func (o *uploadPart) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *uploadPart) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, o)
}

// 初始化分块上传事件
func InitiateMultipartUpload(c *gin.Context) {
	fileName := c.Query("filename")
	list := strings.Split(fileName, ".")
	//fileType := ""
	objectKey := ""
	if len(list) > 1 && (list[len(list)-1] != "") {
		objectKey += strings.ToLower(list[len(list)-1]) + "/"
	}

	objectKey += fileName
	fmt.Printf("%v ==> %v\n", list, objectKey)
	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		//oss.ContentDisposition("attachment;filename=FileName.txt"),
		oss.ContentDisposition("attachment"),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		// oss.ObjectStorageClass(oss.StorageStandard),
	}

	// 步骤1：初始化一个分片上传事件。
	bucket := OSS.BucketConn()
	imur, err := bucket.InitiateMultipartUpload(objectKey, options...)
	imur1 := InitiateMultipartUploadResult(imur)
	//rd := DB.ConnRedis()
	//rd.Set(imur.UploadID, &imur1, 100*time.Second) // 100s后过期
	DB.Set(imur.UploadID, &imur1, 100*time.Second) // 100s后过期
	if err != nil {
		fmt.Println("InitiateMultipartUploadResult Error:", err)
		os.Exit(-1)
	}
	c.JSON(http.StatusOK, gin.H{
		"meta":     models.Meta{0, "InitiateMultipartUpload success"},
		"UploadID": imur.UploadID,
	})
}
