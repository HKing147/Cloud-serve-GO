package main

import (
	"Cloud-serve/src/OSS"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func uploadFile(c *gin.Context) {
	fmt.Print(c)
	//file, _ := c.GetPostForm("file")
	file, _ := c.FormFile("file")
	//fmt.Println("====>", file)

	//OSS.UploadFile(file)
	OSS.MultipartUpload(file)
	savePath := "./" + file.Filename
	filePtr, err := file.Open()
	if err != nil {
		fmt.Println("file error")
		return
	}
	defer filePtr.Close()

	saveFile, err := os.Create(savePath)
	if err != nil {
		fmt.Println("saveFile error")
		return
	}
	defer saveFile.Close()
	var context []byte = make([]byte, 1024)
	for {
		n, err := filePtr.Read(context)
		saveFile.Write(context[:n])
		if err != nil {
			if err == io.EOF {
				return
			}
		}
	}
}
