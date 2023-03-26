package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/OSS"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// 完成分片上传，合并所有分片
func CompleteMultipartUpload(c *gin.Context) {
	UploadID := c.DefaultPostForm("UploadID", "")
	path := c.DefaultPostForm("path", "") // 文件父文件夹路径
	MD5 := c.DefaultPostForm("MD5", "")   // 文件MD5值
	fmt.Println("UploadID:", UploadID)
	bucket := OSS.BucketConn()
	imur := InitiateMultipartUploadResult{}
	parts := []uploadPart{}
	//rd := DB.ConnRedis()
	//rd.Get(UploadID).Scan(&imur)
	//rd.SMembers(UploadID + "_parts").ScanSlice(&parts)
	DB.Get(UploadID).Scan(&imur)
	DB.SMembers(UploadID + "_parts").ScanSlice(&parts)

	parts1 := make([]oss.UploadPart, len(parts))
	for idx, part := range parts {
		parts1[idx] = oss.UploadPart(part)
	}

	cmur, err := bucket.CompleteMultipartUpload(oss.InitiateMultipartUploadResult(imur), parts1, oss.ObjectACL(oss.ACLDefault))

	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	/*
		fmt.Printf("cmur: %v\n", cmur)
		fmt.Printf("cmur.Location: %v\n", cmur.Location) // 文件访问URL
		fmt.Printf("cmur.XMLName.Local: %v\n", cmur.XMLName.Local)
		fmt.Printf("cmur.XMLName.Space: %v\n", cmur.XMLName.Space)
		fmt.Printf("cmur.ETag: %v\n", cmur.ETag) // MD5值
		fmt.Printf("cmur.Key: %v\n", cmur.Key)   // 文件名(路径)
	*/

	fmt.Printf("cmur.Key: %v\n", cmur.Key) // 文件名(路径)
	// 修改文件名为: MD5值+后缀
	srcObject := cmur.Key
	list := strings.Split(srcObject, "/")
	destObject := MD5
	Type := ""
	if len(list) > 1 {
		Type = list[0]
		destObject = Type + "/" + MD5 + "." + Type
	}
	// 复制一份
	fmt.Printf("%v ==> %v\n", srcObject, destObject)
	_, err = bucket.CopyObject(srcObject, destObject)
	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	// 删除原来的
	err = bucket.DeleteObject(srcObject)
	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	//header, _ := bucket.GetObjectMeta(destObject)
	//fmt.Println("文件大小", header.Get("Content-Length"))
	// 向File表中插入文件信息
	/*
		fileUrl: 文件在OSS上的访问链接
		size: 文件大小
		Type: 文件类型
	*/
	size, _ := strconv.ParseInt(c.DefaultPostForm("size", "0"), 10, 64)
	urlDecode, _ := url.QueryUnescape(cmur.Location)
	fmt.Printf("decode: %v  =>> %v", cmur.Location, urlDecode)
	fileUrl := strings.ReplaceAll(urlDecode, srcObject, destObject)
	//var size int64 = 0
	//parentID := 0
	//err = models.InsertFile(fileName, filePath, fileUrl, isFolder, size, Type, &parentID)
	file, err := models.InsertFile(fileUrl, size, Type, MD5)
	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	// 向UserFile表中插入user-file信息
	/*
		userID
		fileID
		fileName
		filePath
	*/
	userID, _ := c.Get("userID")
	fileID := file.ID
	fileName := list[len(list)-1]
	//fmt.Printf("cmur.Key: %v\n", cmur.Key) // 文件名(路径)
	//filePath := path + fileName
	//isFolder := false
	err = models.InsertUserFile(userID.(uint), fileID, fileName, path, false)
	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}

	// User表usedSapce字段更新
	err = models.UpdateUsedSpace(userID.(uint), size)
	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	// 将文件MD5与文件ID写入Redis（永久有效）
	DB.Set(MD5, fileID, 0)

	c.JSON(http.StatusOK, gin.H{
		"meta": models.Meta{0, "CompleteMultipartUpload success"},
	})
}
