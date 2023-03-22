package models

import (
	"Cloud-serve/src/DB"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Share struct {
	gorm.Model
	UserID uint `json:"userID"`
}

type ShareFileList struct {
	UserFileIDList []uint `json:"userfileIDList"`
}

func (o *ShareFileList) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *ShareFileList) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, o)
}

func InsertShare(userID uint, userFileIDList []uint) error {
	// 先创建Share记录
	share := Share{UserID: userID}
	err := db.Create(&share).Error
	if err != nil {
		return err
	}
	// 再向Redis中插入<share_userID_shareID,userFileIDList>
	return DB.Set("share_"+strconv.Itoa(int(userID))+"_"+strconv.Itoa(int(share.ID)), &ShareFileList{userFileIDList}, 30*24*time.Hour).Err()
}

type GetShareListResp struct {
	ShareID   uint      `json:"shareID"`
	FileName  string    `json:"fileName"` // 分享描述(文件名)
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	IsFolder  bool      `json:"isFolder"` // 多个文件一起分享放在一个文件夹里
	CreatedAt time.Time `json:"createdTime"`
	UpdatedAt time.Time `json:"updatedTime"`
}

// 获取用户的分享列表
func GetShareList(userID uint) ([]GetShareListResp, error) {
	// 先获取到所有的shareID
	shareList := []Share{}
	err := db.Model(&Share{}).Where("user_id = ?", userID).Scan(&shareList).Error
	if err != nil {
		return nil, err
	}
	// 再通过shareID获取shareList
	res := make([]GetShareListResp, len(shareList))
	for i, share := range shareList {
		//share := Share{}
		//err = DB.Get("share_" + strconv.Itoa(int(userID)) + "_" + strconv.Itoa(int(shareID))).Scan(&share)
		//if err != nil {
		//	return nil, err
		//}
		// 获取文件列表
		fileList, err := GetShareByID(share.ID)
		if err != nil {
			return nil, err
		}
		res[i].ShareID = share.ID
		res[i].CreatedAt = share.CreatedAt
		res[i].UpdatedAt = share.UpdatedAt
		// 判断是否只有一个文件（夹）
		if len(fileList) == 1 {
			res[i].Type = fileList[0].Type
			res[i].FileName = fileList[0].FileName
			res[i].IsFolder = fileList[0].IsFolder
			res[i].Size = fileList[0].Size
		} else if len(fileList) > 1 { // 大于一个文件
			res[i].Type = "folder" // 变成文件夹
			res[i].IsFolder = true
			res[i].FileName = fmt.Sprintf("%v 等 %v 个文件", fileList[0].FileName, len(fileList))
		}
	}
	return res, err
}

// 通过shareID获得分享的文件
func GetShareByID(shareID uint) ([]SelectFilesByUserIDAndPathResp, error) {
	// 先获得分享者ID
	var userID uint
	err := db.Model(&Share{}).Where("id = ?", shareID).Select("user_id").Scan(&userID).Error
	if err != nil {
		return nil, err
	}
	// 然后再获取userFileIDList
	shareFileList := ShareFileList{}
	err = DB.Get("share_" + strconv.Itoa(int(userID)) + "_" + strconv.Itoa(int(shareID))).Scan(&shareFileList)
	if err != nil {
		return nil, err
	}
	// 最后在获取具体的文件列表
	fileList := []SelectFilesByUserIDAndPathResp{}
	err = db.Model(&UserFile{}).Where("`user_files`.user_id = ? and `user_files`.id in ?", userID, shareFileList.UserFileIDList).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList).Error
	return fileList, err
}
