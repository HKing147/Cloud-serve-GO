package main

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/OSS"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"os"
)

// 完成分片上传，合并所有分片
func CompleteMultipartUpload(c *gin.Context) {
	UploadID := c.DefaultPostForm("UploadID", "")
	fmt.Println("UploadID:", UploadID)
	bucket := OSS.BucketConn()
	rd := DB.ConnRedis()
	imur := InitiateMultipartUploadResult{}
	rd.Get(UploadID).Scan(&imur)
	parts := []uploadPart{}
	rd.SMembers(UploadID + "_parts").ScanSlice(&parts)
	parts1 := make([]oss.UploadPart, len(parts))
	for idx, part := range parts {
		parts1[idx] = oss.UploadPart(part)
	}

	cmur, err := bucket.CompleteMultipartUpload(oss.InitiateMultipartUploadResult(imur), parts1, oss.ObjectACL(oss.ACLDefault))

	if err != nil {
		fmt.Println("CompleteMultipartUpload Error:", err)
		os.Exit(-1)
	}
	fmt.Printf("cmur: %v\n", cmur)
}
