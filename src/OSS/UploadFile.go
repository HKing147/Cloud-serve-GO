package OSS

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"os"
)

func UploadFile(file *multipart.FileHeader) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", "LTAI5tEanogqiF8wvgGm1H7k", "y9u4T00B6WsqLkPkaGu5uketQrJ4ph")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket("file-bucket001")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写本地文件的完整路径，例如D:\\localpath\\examplefile.txt。
	//fd, err := os.Open("D:\\localpath\\examplefile.txt")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//defer fd.Close()
	fmt.Println("##########################################")
	filePtr, err := file.Open()

	if err != nil {
		fmt.Println("OSS error...")
	}
	defer filePtr.Close()
	// 将文件流上传至exampledir目录下的exampleobject.txt文件。
	//content := make([]byte, 1024*38)
	//filePtr.Read(content)
	//content, _ := ioutil.ReadAll(filePtr)

	//err = bucket.PutObject("exampledir/"+file.Filename, bytes.NewReader(content))
	err = bucket.PutObject("exampledir/"+file.Filename, filePtr) // ！！！
	//err = bucket.UploadFile("exampledir/"+file.Filename, file, 100*1024)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

}
