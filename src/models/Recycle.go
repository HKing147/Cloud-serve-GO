package models

import (
	"fmt"
	"gorm.io/gorm"
)

// 回收站
type Recycle struct {
	gorm.Model
	UserID     uint     `json:"userID"`
	UserFileID uint     `json:"userFileID"` // UserFile表ID
	userFile   UserFile `json:"userFile" gorm:"Foreignkey:UserFileID"`
}

func InsertRecycle(userID uint, userFileID uint) error {
	return db.Create(&Recycle{UserID: userID, UserFileID: userFileID}).Error
}

// 获取用户删除的文件列表
func GetRecycle(userID uint) ([]SelectFilesByUserIDAndPathResp, error) {
	list := []SelectFilesByUserIDAndPathResp{}
	res := db.Model(&Recycle{}).Joins(", `user_files`").Joins(", `files`").Where("`recycles`.user_id = ? and `recycles`.user_file_id = `user_files`.id and `user_files`.file_id = `files`.id", userID).Select("`user_files`.*, `files`.*, `user_files`.id as id").Scan(&list)
	// SELECT `user_files`.*, `files`.*, `user_files`.id as id FROM `recycles`, `user_files`, `files` WHERE (`recycles`.user_id = 1 and `recycles`.user_file_id = `user_files`.id and `user_files`.file_id = `files`.id) AND `recycles`.`deleted_at` IS NULL;resumeFiles
	return list, res.Error
}

// 取消删除文件
//func ResumeRecycle(userID uint, userFileIDList []uint) error {
//	// 删除这些记录
//	return db.Model(&Recycle{}).Where("user_id = ? and user_file_id in ?", userID, userFileIDList).Delete(&Recycle{}).Error
//}

// 取消删除文件
func ResumeRecycle(userID uint, userFileID uint) error {
	fmt.Printf("ResumeRecycle ==> userID: %v, userFileID: %v\n", userID, userFileID)
	// 删除这些记录
	return db.Model(&Recycle{}).Where("user_id = ? and user_file_id = ?", userID, userFileID).Delete(&Recycle{}).Error
	// return db.Model(&Recycle{}).Delete("user_id = ? and user_file_id = ?", userID, userFileID).Error
}
