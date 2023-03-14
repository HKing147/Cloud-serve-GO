package OSS

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type UploadPart oss.UploadPart

func (part *UploadPart) MarshalBinary() (data []byte, err error) {
	return json.Marshal(part)
}

func (part *UploadPart) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, part)
}

func MultipartUpload(file *multipart.FileHeader) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", "LTAI5tEanogqiF8wvgGm1H7k", "y9u4T00B6WsqLkPkaGu5uketQrJ4ph")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 填写存储空间名称。
	bucketName := "file-bucket001"
	// 填写Object完整路径。Object完整路径中不能包含Bucket名称。
	//objectName := "abc.png"

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		// oss.ContentDisposition("attachment;filename=FileName.txt"),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		// oss.ObjectStorageClass(oss.StorageStandard),
	}

	// 步骤1：初始化一个分片上传事件。
	imur, err := bucket.InitiateMultipartUpload(file.Filename, options...)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	filePtr, err := file.Open()

	if err != nil {
		fmt.Println("OSS error...")
	}
	defer filePtr.Close()

	var chunkSize int64 = 1024 * 1024 // 1MB

	// imur file idx

	// 步骤2：上传分片。
	var parts []oss.UploadPart
	var currentChunk int64 = 0
	var idx int = 1
	for {
		if currentChunk > file.Size {
			break
		}
		filePtr.Seek(int64(currentChunk), os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		//part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		var part oss.UploadPart
		if file.Size-currentChunk >= chunkSize {
			part, err = bucket.UploadPart(imur, filePtr, chunkSize, idx)
		} else {
			part, err = bucket.UploadPart(imur, filePtr, file.Size-currentChunk, idx)
		}
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		fmt.Printf("part.ETag: %v  part.PartNumber: %v  part.XMLName: %v\n", part.ETag, part.PartNumber, part.XMLName)
		parts = append(parts, part)
		fmt.Println("Completed", idx)
		currentChunk += chunkSize
		idx++
	}

	// 指定Object的读写权限为私有，默认为继承Bucket的读写权限。
	//objectAcl := oss.ObjectACL(oss.ACLPrivate)

	// 步骤3：完成分片上传。
	//cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	fmt.Printf("parts: %v\n", parts)
	cmur, err := bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("cmur:", cmur)
}
