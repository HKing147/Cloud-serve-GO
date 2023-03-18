package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// UserFile
type UserFile struct {
	gorm.Model
	FileID    uint
	File      File `json:"file" gorm:"Foreignkey:FileID"` // 文件
	UserID    uint
	User      User   `json:"user" gorm:"Foreignkey:UserID"` // 文件所属用户
	FileName  string `json:"fileName"`
	FilePath  string `json:"filePath"`
	IsFolder  bool   `json:"isFolder"`  // 是否是非文件夹
	IsCollect bool   `json:"isCollect"` // 用户是否收藏了该文件
	IsShare   bool   `json:"isShare"`   // 用户是否分享了该文件
}

func GetFileList(c *gin.Context) {
	userID, _ := c.Get("userID")
	path := c.Query("path")
	fileList := SelectFilesByUserIDAndPath(userID.(uint), path)
	c.JSON(http.StatusOK, gin.H{"meta": Meta{0, "success"}, "fileList": fileList})
}

func InsertUserFile(userID uint, fileID uint, fileName string, filePath string, isFolder bool) error {
	userFile := UserFile{
		UserID:   userID,
		FileID:   fileID,
		FileName: fileName,
		FilePath: filePath,
		IsFolder: isFolder,
	}
	res := db.Create(&userFile)
	return res.Error
}

type SelectFilesByUserIDAndPathResp struct {
	ID          uint      `json:"id"`
	FileName    string    `json:"fileName"`
	FileUrl     string    `json:"fileUrl"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"`
	NotFolder   bool      `json:"notFolder"`
	IsCollected bool      `json:"isCollected"`
	IsShare     bool      `json:"isShare"`
	UpdatedAt   time.Time `json:"updatedTime"`
}

func SelectFilesByUserIDAndPath(userID uint, path string) []SelectFilesByUserIDAndPathResp {
	//fileList := []SelectFilesByUserIDAndPathResp{}
	res := []SelectFilesByUserIDAndPathResp{}
	//db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_path like ?", userID, path+"%").Joins("File").Select("`user_files`.*, File.*").Scan(&res)
	db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_path = ?", userID, path).Joins("File").Select("`user_files`.*, File.*").Scan(&res)
	return res
}
