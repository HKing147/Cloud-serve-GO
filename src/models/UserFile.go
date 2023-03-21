package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
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

// 只删除当前文件（Recycle表和UserFile表）
func DeleteOneFile(userID uint, file UserFile) error {
	// Recycle表中添加一条记录
	err := InsertRecycle(userID, file.ID)
	if err != nil {
		return err
	}
	// UserFile表中删除这条记录: DeletedAt置为NULL
	return db.Delete(&file).Error
}

// 递归进入文件夹删除
func DeleteFileDown(userID uint, folder UserFile) error {
	// 获取文件夹下的文件列表
	fileList := []UserFile{}
	err := db.Model(&UserFile{}).Where("user_id = ? and file_path = ?", userID, folder.FilePath+folder.FileName+"/").Scan(&fileList).Error
	if err != nil {
		return err
	}

	for _, file := range fileList {
		// 判断是否是文件夹
		if file.IsFolder { // 是文件夹，递归进入删除
			err := DeleteFileDown(userID, file)
			if err != nil {
				return err
			}
		}
		// 删除文件（不用向Recycle表中插入）
		err = db.Delete(&file).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// 只有第一层的文件才会向Recycle表中插入一条记录
func DeleteFiles(userID uint, userFileIDList []uint, path string) error {
	//return db.Model(&UserFile{}).Delete("user_id = ? and file_id in ? file_path = ?", userID, fileIDList, path).Error
	// 递归删除
	for _, userFileID := range userFileIDList {
		// 先看它是否是文件夹
		file := UserFile{}
		err := db.Model(&UserFile{}).Where("id = ? and user_id = ? and file_path = ?", userFileID, userID, path).Scan(&file).Error
		if err != nil {
			return err
		}
		//db.Model(&UserFile{}).Delete("file_id = ? and user_id = ? and file_path = ?", fileID, userID, path)
		if file.IsFolder { // 是文件夹
			// 递归进入删除
			//err = DeleteFileDown(userID, []UserFile{file}, path+file.FileName+"/")
			err = DeleteFileDown(userID, file)
			if err != nil {
				return err
			}
		}
		// 删除该文件, 不管是否是文件夹
		err = DeleteOneFile(userID, file) // 需要向Recycle表中插入一条记录
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取用户删除的文件列表
//func GetRecycle(userID uint) ([]SelectFilesByUserIDAndPathResp, error) {
//	list := []SelectFilesByUserIDAndPathResp{}
//	// deleted_at 字段不为NULL ==> 该记录之前被删除了
//	res := db.Unscoped().Model(&UserFile{}).Where("`user_files`.user_id = ? and `user_files`.deleted_at is not null", userID).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&list)
//	return list, res.Error
//}

// 只单单恢复当前文件
func ResumeOneFile(userID uint, file UserFile) error {
	fmt.Println("恢复文件", file.FileName, file.FilePath, file.IsFolder)
	//err := ResumeRecycle(userID, []uint{file.ID})
	err := ResumeRecycle(userID, file.ID)
	if err != nil {
		return err
	}
	// 再将UserFile中DeletedAt字段置为NULL
	return db.Unscoped().Model(&file).Update("deleted_at", nil).Error
}

// 向上恢复文件夹(一定要是文件夹！！！)
func ResumeFolderUP(userID uint, son UserFile) error {
	if son.FilePath == "/" { // 已经恢复到根文件夹了，直接退出
		return nil
	}
	strList := strings.Split(son.FilePath, "/")
	fmt.Printf("len: %v strList: %v\n", len(strList), strList)
	// 获取父文件夹名
	parentName := strList[len(strList)-2]
	// 获取父文件夹所在文件夹路径
	parentPath := son.FilePath[:len(son.FilePath)-len(parentName)-1]
	// 查出父文件夹
	parent := UserFile{}
	err := db.Unscoped().Model(&UserFile{}).Where("file_name = ? and file_path = ? and is_folder = ?", parentName, parentPath, true).Scan(&parent).Error
	if err != nil {
		return err
	}
	// 先检查父文件夹此时是否存在（可能被删除了，但是在其它子文件恢复时将其恢复了；也有可能没被删除过）
	// 防止重复恢复，提升性能
	if !parent.DeletedAt.Valid { // DeletedAt为空，存在，直接return，不用继续往上恢复
		return nil
	}
	// 父文件夹不存在，将其恢复
	//err = ResumeOneFile(userID, parent)// 错误
	err = db.Unscoped().Model(&parent).Update("deleted_at", nil).Error // 不用删除Recycle表的记录(不是第一层)
	if err != nil {
		return err
	}
	/*
		// 先将Recycle表中记录删除
		err = ResumeRecycle(userID, []uint{parent.ID})
		if err != nil {
			return err
		}
		// 再将UserFile中DeletedAt字段置为NULL
		err = db.Unscoped().Model(&parent).Update("deleted_at", nil).Error
		if err != nil {
			return err
		}
	*/
	// 最后再继续向上恢复
	return ResumeFolderUP(userID, parent)
}

// 向下递归恢复文件（不用再向上了）
func ResumeFilesDown(userID uint, parent UserFile) error {
	// 先查询出它的子文件
	path := parent.FilePath + parent.FileName + "/"
	fileList := []UserFile{}
	err := db.Unscoped().Model(&UserFile{}).Where("user_id = ? and file_path = ?", userID, path).Scan(&fileList).Error
	if err != nil {
		return err
	}
	for _, file := range fileList {
		if file.IsFolder { // 是文件夹继续就递归
			err = ResumeFilesDown(userID, file)
			if err != nil {
				return err
			}
		}
		// 恢复当前文件/文件夹
		//err = ResumeOneFile(userID, file)// 错误
		err = db.Unscoped().Model(&file).Update("deleted_at", nil).Error // 不用删除Recycle表的记录(不是第一层)
		if err != nil {
			return err
		}
	}
	return nil
}

// 恢复被删除的文件(只有第一层才会将UserFile表中的deleted_at置为NULL)
func ResumeFiles(userID uint, userFileIDList []uint) error {
	/**
	简单恢复
	*/
	//// 先将Recycle表中的记录删除
	//err := ResumeRecycle(userID, userFileIDList)
	//if err != nil {
	//	fmt.Println("提前退出")
	//	return err
	//}
	//// 再将UserFile表中的deleted_at字段置为NULL
	//// 要加Unscoped()！！！
	//return db.Unscoped().Model(&UserFile{}).Where("user_id = ? and id in ?", userID, userFileIDList).Update("deleted_at", nil).Error

	/**
	TODO: file如果是文件夹还要做到：
		1.向上恢复包含它的文件夹
		2.向下递归恢复其子文件
	*/
	// 先判断是否是文件夹
	for _, userFileID := range userFileIDList {
		file := UserFile{}
		db.Unscoped().Model(&UserFile{}).Where("user_id = ? and id = ?", userID, userFileID).Scan(&file)

		// 1.向上恢复包含它的父文件夹(文件和文件夹都要)
		err := ResumeFolderUP(userID, file)
		if err != nil {
			return err
		}
		if file.IsFolder { // 是文件夹
			// 2.向下递归恢复其子文件（是由文件夹才要递归下去）
			err = ResumeFilesDown(userID, file)
			if err != nil {
				return err
			}
		}
		// 不管是文件夹还是文件，都要将Recycle表中记录删除，(只有第一层)将UserFile表中的deleted_at置为NULL
		err = ResumeOneFile(userID, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// 查询文件
func SearchFile(userID uint, fileName string) ([]SelectFilesByUserIDAndPathResp, error) {
	fileList := []SelectFilesByUserIDAndPathResp{}
	res := db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_name like ?", userID, "%"+fileName+"%").Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList)
	return fileList, res.Error
}

// 递归(只)更新path
func updatePathDown(userID uint, oldPath string, newPath string) error {
	// 查询出它的子文件
	fileList := []UserFile{}
	err := db.Model(&UserFile{}).Where("user_id = ? and file_path = ?", userID, oldPath).Scan(&fileList).Error
	if err != nil {
		return err
	}
	for _, file := range fileList {
		err = db.Model(&file).Update("file_path", newPath).Error
		if err != nil {
			return err
		}
		// 是文件夹继续递归
		if file.IsFolder {
			err = updatePathDown(userID, oldPath+file.FileName+"/", newPath+file.FileName+"/")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 重命名文件
func RenameFile(userID uint, userFileID uint, newFileName string) error {
	file := UserFile{}
	err := db.Model(&UserFile{}).Where("user_id = ? and id = ?", userID, userFileID).Scan(&file).Error
	if err != nil {
		return err
	}
	// 先存下旧文件名
	oldFileName := file.FileName
	// 修改当前文件名
	err = db.Model(&file).Update("file_name", newFileName).Error
	if err != nil {
		return err
	}
	// 如果修改的是文件夹，则还要将其子文件的path也修改
	if file.IsFolder {
		err = updatePathDown(userID, file.FilePath+oldFileName+"/", file.FilePath+newFileName+"/")
		if err != nil {
			return err
		}
	}
	return nil
}
