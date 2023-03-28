package OSS

import (
	"fmt"
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

	err = bucket.PutObject(objName, filePtr) // ！！！
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return "https://" + BucketName + "." + Endpoint + "/" + objName, nil
}
