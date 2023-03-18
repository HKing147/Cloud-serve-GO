package OSS

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

func BucketConn() *oss.Bucket {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", "LTAI5tRnVo7L548yNZ2bRJxv", "X848KXtK7imYFSAjcCQUK0W6TH4cTR")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 填写存储空间名称。
	bucketName := "file-bucket001"

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return bucket
}
