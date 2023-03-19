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
	fileList, err := SelectFilesByUserIDAndPath(userID.(uint), path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "error"}})
		return
	}
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
	ID        uint      `json:"id"`
	FileName  string    `json:"fileName"`
	FileUrl   string    `json:"fileUrl"`
	FilePath  string    `json:"filePath"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	NotFolder bool      `json:"notFolder"`
	IsCollect bool      `json:"isCollect"`
	IsShare   bool      `json:"isShare"`
	UpdatedAt time.Time `json:"updatedTime"`
}

func SelectFilesByUserIDAndPath(userID uint, path string) ([]SelectFilesByUserIDAndPathResp, error) {
	//fileList := []SelectFilesByUserIDAndPathResp{}
	fileList := []SelectFilesByUserIDAndPathResp{}
	//db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_path like ?", userID, path+"%").Joins("File").Select("`user_files`.*, File.*").Scan(&res)
	res := db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_path = ?", userID, path).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList)
	return fileList, res.Error
}

func CollectedFile(userID uint, fileIDList []uint) error {
	//return db.Where("user_id = ? and file_id in ?", userID, fileIDList).Exec("is_collected = !is_collected").Error
	//return db.Model(&UserFile{}).Exec("set is_collect = if(is_collect,0,1) where user_id = ? and id in ?", userID, fileIDList).Error
	return db.Exec("update `user_files` set is_collect = if(is_collect,0,1) where user_id = ? and id in ?", userID, fileIDList).Error
}

func GetCollectedList(userID uint) ([]SelectFilesByUserIDAndPathResp, error) {
	collectedList := []SelectFilesByUserIDAndPathResp{}
	res := db.Model(&UserFile{}).Where("`user_files`.user_id = ? and is_collect = true", userID).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&collectedList)
	return collectedList, res.Error
}

// 获取相册，文件类型为: png, jpg, jpeg, webp, gif, ico
func GetAlbum(userID uint) ([]SelectFilesByUserIDAndPathResp, error) {
	typeList := []string{"png", "jpg", "jpeg", "webp", "gif", "ico"} // 图片类型列表
	collectedList := []SelectFilesByUserIDAndPathResp{}
	res := db.Model(&UserFile{}).Where("`user_files`.user_id = ? and type in ?", userID, typeList).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&collectedList)
	return collectedList, res.Error
}

func DeleteFile(userID uint, folderID uint, path string) error {
	// 获取文件夹下的文件列表
	//type T struct {
	//	FileID uint `json:"fileID"`
	//}
	fileList := []UserFile{}
	res := db.Model(&UserFile{}).Where("user_id = ? and file_path = ?", userID, path).Scan(&fileList)
	if res.Error != nil {
		return res.Error
	}
	var fileIDList []uint
	fileIDList = make([]uint, len(fileList))
	for i, item := range fileList {
		fileIDList[i] = item.FileID
	}
	err := DeleteFiles(userID, fileIDList, path)
	return err
}

func DeleteFiles(userID uint, fileIDList []uint, path string) error {
	/**
	TODO: file如果是文件夹还要做到递归删除
	*/
	//return db.Model(&UserFile{}).Delete("user_id = ? and file_id in ? file_path = ?", userID, fileIDList, path).Error
	// 递归删除
	for _, fileID := range fileIDList {
		// 先看它是否是文件夹
		item := UserFile{}
		err := db.Model(&UserFile{}).Where("file_id = ? and user_id = ? and file_path = ?", fileID, userID, path).Scan(&item)
		if err.Error != nil {
			return err.Error
		}
		//db.Model(&UserFile{}).Delete("file_id = ? and user_id = ? and file_path = ?", fileID, userID, path)
		if item.IsFolder { // 是文件夹
			// 递归进入删除
			err := DeleteFile(userID, fileID, path+item.FileName+"/")
			if err != nil {
				return err
			}
		}
		// 删除该文件, 不管是否是文件夹
		err = db.Delete(&item)
		if err.Error != nil {
			return err.Error
		}
	}
	return nil
}

// 获取用户删除的文件列表
func GetRecycle(userID uint) ([]SelectFilesByUserIDAndPathResp, error) {
	list := []SelectFilesByUserIDAndPathResp{}
	// deleted_at 字段不为NULL ==> 该记录之前被删除了
	res := db.Unscoped().Model(&UserFile{}).Where("`user_files`.user_id = ? and `user_files`.deleted_at is not null", userID).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&list)
	return list, res.Error
}
