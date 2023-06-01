package OSS

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var (
	Endpoint        string = "oss-cn-hangzhou.aliyuncs.com"
	AccessKeyID     string = "LTAI5tRnVo7L548yNZ2bRJxv"
	AccessKeySecret string = "X848KXtK7imYFSAjcCQUK0W6TH4cTR"
	BucketName      string = "file-bucket001"
)

func Client() *oss.Client {
	// 创建OSSClient实例。
	client, err := oss.New(Endpoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	return client
}

func BucketConn() *oss.Bucket {
	client := Client()
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return bucket
}
