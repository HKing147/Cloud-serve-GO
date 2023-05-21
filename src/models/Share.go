package models

import (
	"Cloud-serve/src/DB"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Share struct {
	gorm.Model
	UserID         uint           `json:"userID"`
	Password       string         `json:"password"`
	ExpirationTime gorm.DeletedAt `json:"expirationTime" gorm:"default:null"` // 过期时间
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

func InsertShare(userID uint, userFileIDList []uint, shareMethod bool, shareDuration int) (Share, error) {
	share := Share{UserID: userID}
	if shareMethod { // 需要密码
		password := ""
		for i := 0; i < 6; i++ {
			password += string('a' + rand.Intn(26))
		}
		share.Password = password
	}
	if shareDuration != 0 { // 有有效期
		share.ExpirationTime = gorm.DeletedAt{time.Now().AddDate(0, 0, shareDuration), true} // YY MM DD, NullTime.Valid = true ==> 非空
	}
	// 先创建Share记录
	err := db.Create(&share).Error
	if err != nil {
		return Share{}, err
	}
	// 再向Redis中插入<share_userID_shareID,userFileIDList>
	err = DB.Set("share_"+strconv.Itoa(int(userID))+"_"+strconv.Itoa(int(share.ID)), &ShareFileList{userFileIDList}, 30*24*time.Hour).Err()
	return share, err
}

func DeleteShares(userID uint, shareIDList []uint) error {
	// 从Share表中删除
	err := db.Where("user_id = ? and id in ?", userID, shareIDList).Delete(&Share{}).Error
	if err != nil {
		return err
	}
	// 从redis中删除
	keys := make([]string, len(shareIDList))
	for i, shareID := range shareIDList {
		key := fmt.Sprintf("share_%v_%v", userID, shareID)
		keys[i] = key
	}
	return DB.Del(keys...).Err()
}

type GetShareListResp struct {
	//ShareID   uint      `json:"shareID"`
	ID        uint      `json:"id"`
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
	err := db.Model(&Share{}).Where("user_id = ? and (expiration_time is null or expiration_time > ?)", userID, time.Now()).Scan(&shareList).Error
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
		_, fileList, err := GetShareByID(share.ID)
		if err != nil {
			return nil, err
		}
		//res[i].ShareID = share.ID
		res[i].ID = share.ID
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
func GetShareByID(shareID uint) (Share, []SelectFilesByUserIDAndPathResp, error) {
	// 先获得分享者ID
	//var userID uint
	share := Share{}
	err := db.Model(&Share{}).Where("id = ? and (expiration_time is null or expiration_time > ?)", shareID, time.Now()).Scan(&share).Error
	if err != nil {
		return share, nil, err
	}
	// 然后再获取userFileIDList
	shareFileList := ShareFileList{}
	err = DB.Get("share_" + strconv.Itoa(int(share.UserID)) + "_" + strconv.Itoa(int(shareID))).Scan(&shareFileList)
	if err != nil {
		return share, nil, err
	}
	// 最后在获取具体的文件列表
	fileList := []SelectFilesByUserIDAndPathResp{}
	err = db.Model(&UserFile{}).Order("is_folder desc").Where("`user_files`.user_id = ? and `user_files`.id in ?", share.UserID, shareFileList.UserFileIDList).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList).Error
	return share, fileList, err
}

// 通过shareUrl获得分享的文件
func GetShareByShareUrl(shareUrl string, sortMethod string) (Share, []SelectFilesByUserIDAndPathResp, error) {
	if strings.HasPrefix(sortMethod, "updated_at") {
		sortMethod = "`user_files`." + sortMethod
	}

	list := strings.Split(shareUrl, "_")
	fmt.Println(list)
	if len(list) != 3 || list[0] != "share" {
		return Share{}, nil, errors.New("error")
	}
	userID, err := strconv.Atoi(list[1])
	if err != nil {
		return Share{}, nil, err
	}
	shareID, err := strconv.Atoi(list[2])
	if err != nil {
		return Share{}, nil, err
	}
	// 获取Share信息（密码，过期时间）
	shareInfo := Share{}
	err = db.Model(&Share{}).Where("id = ? and user_id = ? and (expiration_time is null or expiration_time > ?)", shareID, userID, time.Now()).Scan(&shareInfo).Error
	if err != nil {
		return Share{}, nil, err
	}
	shareFileList := ShareFileList{}
	err = DB.Get(shareUrl).Scan(&shareFileList)
	if err != nil {
		return Share{}, nil, err
	}
	// 获取具体的文件列表
	fileList := []SelectFilesByUserIDAndPathResp{}
	err = db.Model(&UserFile{}).Order("is_folder desc").Order(sortMethod).Where("`user_files`.user_id = ? and `user_files`.id in ?", userID, shareFileList.UserFileIDList).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList).Error
	return shareInfo, fileList, err
}

func GetShareCnt() (shareCnt int64, err error) {
	err = db.Model(&Share{}).Count(&shareCnt).Error
	return
}

// 统计一周内每天的分享数
func GetShareCntByDay(t time.Time) (shareCnt int64, err error) {
	err = db.Model(&Share{}).Where("date(created_at) = ?", t).Count(&shareCnt).Error
	return
}
