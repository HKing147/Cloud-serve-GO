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
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	//client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", "LTAI5tRnVo7L548yNZ2bRJxv", "X848KXtK7imYFSAjcCQUK0W6TH4cTR")
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
