package OSS

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
)

func UploadFile(file *multipart.FileHeader, objName string) (string, error) {
	bucket := BucketConn()
	filePtr, err := file.Open()
	if err != nil {
		fmt.Println("OSS error...")
		return "", err
	}
	defer filePtr.Close()
	options := []oss.Option{
		// 指定该Object被下载时的名称。
		//oss.ContentDisposition("attachment;filename=FileName.txt"),
		oss.ContentDisposition("attachment"),
	}
	err = bucket.PutObject(objName, filePtr, options...) // ！！！
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return "https://" + BucketName + "." + Endpoint + "/" + objName, nil
}
