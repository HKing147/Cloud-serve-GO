package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// File
type File struct {
	gorm.Model
	DownCount int64 `json:"downCount"` // 下载次数
	//FileName  string `json:"fileName"`  // 文件名
	//FilePath  string `json:"filePath"`  // 文件路径（相对于当前用户的）
	FileURL string `json:"fileUrl"` // 文件URL
	//ID         int64  `json:"id"`        // 文件ID
	NotFolder bool `json:"notFolder"` // 是否是文件夹
	//ModifyTime int64  `json:"modifyTime"` // 修改时间
	//ParentID int64  `json:"parentID"` // 父文件夹ID
	Size int64  `json:"size"` // 文件大小
	Type string `json:"type"` // 文件类型
	MD5  string `json:"md5"`  // 文件md5值
	//UploadTime int64  `json:"uploadTime"` // 上传时间
	//ParentID *int
	//Parent   *File `json:"parent"` // 父文件夹
}

// func InsertFile(fileName string, filePath string, fileUrl string, isFolder bool, size int64, Type string) error {
func InsertFile(fileUrl string, size int64, Type string, MD5 string) (*File, error) {
	file := File{
		//FileName: fileName,
		//FilePath: filePath,
		FileURL: fileUrl,
		//IsFolder: isFolder,
		Size: size,
		Type: Type,
		MD5:  MD5,
		//ParentID: parentID,
	}
	res := db.Create(&file)
	return &file, res.Error
}

func GetFileByID(fileID int) (File, error) {
	file := File{}
	err := db.Model(&File{}).Where("id = ?", fileID).Scan(&file).Error
	return file, err
}

func AddDownCount(fileID uint) error {
	return db.Model(&File{}).Where("id = ?", fileID).Update("down_count", gorm.Expr("down_count + 1")).Error
}

type T struct {
	Name  string `gorm:"name" json:"name"`
	Value uint   `gorm:"value" json:"value"`
}

func GetFileCategory() ([]T, error) {
	res := []T{}
	err := db.Model(&File{}).Group("type").Select("type as name, count(*) as value").Order("value desc").Scan(&res).Error
	fmt.Println(res)
	return res, err
}

// 统计一周内每天的文件上传数
func GetUploadCntByDay(t time.Time) (uploadCnt int64, err error) {
	err = db.Model(&File{}).Where("date(created_at) = ? and type != 'folder'", t).Count(&uploadCnt).Error
	return
}
