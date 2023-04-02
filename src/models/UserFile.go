package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"sync"
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
	sortMethod := c.Query("sortMethod")
	fileList, err := SelectFilesByUserIDAndPath(userID.(uint), path, sortMethod)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": Meta{0, "success"}, "fileList": fileList})
}

func GetFileListByFolderID(folderID uint, sortMethod string) ([]SelectFilesByUserIDAndPathResp, error) {
	folder := UserFile{}
	db.Model(&UserFile{}).Where("id = ?", folderID).Scan(&folder)
	return SelectFilesByUserIDAndPath(folder.UserID, folder.FilePath+folder.FileName+"/", sortMethod)
}

func AddUserFile(userID uint, fileID uint, fileName string, filePath string, isFolder bool) (UserFile, error) {
	userFile := UserFile{
		UserID:   userID,
		FileID:   fileID,
		FileName: fileName,
		FilePath: filePath,
		IsFolder: isFolder,
	}
	err := db.Create(&userFile).Error
	if err != nil {
		return userFile, err
	}
	var size int64
	err = db.Model(&File{}).Where("id = ?", fileID).Select("size").Scan(&size).Error
	if err != nil {
		return userFile, err
	}
	// User表usedSapce字段更新
	return userFile, UpdateUsedSpace(userID, size)
}

// 保证同一目录下不存在同名文件
func InsertUserFile(userID uint, fileID uint, fileName string, filePath string, isFolder bool) (userFile UserFile, err error) {
	fileName_ := fileName
	fileType := ""
	if !isFolder { // 不是文件夹（需要抠出文件后缀）
		list := strings.Split(fileName, ".")
		fileName_ = fileName[:len(fileName)-len(list[len(list)-1])-1]
		fileType = "." + list[len(list)-1]
	}
	i := 0

	for {
		fileName := fileName_
		if i != 0 {
			fileName += fmt.Sprintf("(%v)", i)
		}
		// 查询之前目录下是否已经存在同名文件
		var tmp uint
		db.Model(&UserFile{}).Where("user_id = ? and file_path = ? and file_name = ?", userID, filePath, fileName+fileType).Select("id").First(&tmp)
		if tmp == 0 { // 不存在同名文件
			log.Printf("文件：%v可以插入（没有同名文件）\n", fileName+fileType)
			// 插入新文件
			userFile, err = AddUserFile(userID, fileID, fileName+fileType, filePath, isFolder)
			if err != nil {
				return userFile, err
			}
			break
		}
		i++
	}
	return userFile, err
}

/*
	func InsertUserFile(userID uint, fileID uint, fileName string, filePath string, isFolder bool) error {
		userFile := UserFile{
			UserID:   userID,
			FileID:   fileID,
			FileName: fileName,
			FilePath: filePath,
			IsFolder: isFolder,
		}
		err := db.Create(&userFile).Error
		if err != nil {
			return err
		}
		var size int64
		err = db.Model(&File{}).Where("id = ?", fileID).Select("size").Scan(&size).Error
		if err != nil {
			return err
		}
		// User表usedSapce字段更新
		return UpdateUsedSpace(userID, size)
	}
*/

type SelectFilesByUserIDAndPathResp struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"fileName"`
	FileUrl   string    `json:"fileUrl"`
	FilePath  string    `json:"filePath"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	IsFolder  bool      `json:"isFolder"`
	IsCollect bool      `json:"isCollect"`
	IsShare   bool      `json:"isShare"`
	CreatedAt time.Time `json:"createdTime"`
	UpdatedAt time.Time `json:"updatedTime"`
}

func SelectFilesByUserIDAndPath(userID uint, path string, sortMethod string) ([]SelectFilesByUserIDAndPathResp, error) {
	//fileList := []SelectFilesByUserIDAndPathResp{}
	fileList := []SelectFilesByUserIDAndPathResp{}
	//db.Model(&UserFile{}).Where("`user_files`.user_id = ? and file_path like ?", userID, path+"%").Joins("File").Select("`user_files`.*, File.*").Scan(&res)
	if strings.HasPrefix(sortMethod, "updated_at") {
		sortMethod = "`user_files`." + sortMethod
	}
	res := db.Model(&UserFile{}).Order("is_folder desc").Order(sortMethod).Where("`user_files`.user_id = ? and file_path = ?", userID, path).Joins("File").Select("`user_files`.*, File.*, `user_files`.id as id").Scan(&fileList)
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
	err = db.Delete(&file).Error
	if err != nil {
		return err
	}
	// 更新usedSpace
	var size int64
	db.Model(&File{}).Where("id = ?", file.FileID).Select("size").Scan(&size)
	return UpdateUsedSpace(userID, -size)
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
		// 更新usedSpace
		var size int64
		db.Model(&File{}).Where("id = ?", file.FileID).Select("size").Scan(&size)
		err = UpdateUsedSpace(userID, -size)
		if err != nil {
			return err
		}
	}
	return nil
}

// 只有第一层的文件才会向Recycle表中插入一条记录
func DeleteFiles(userID uint, userFileIDList []uint, path string) error {
	//return db.Model(&UserFile{}).Delete("user_id = ? and file_id in ? file_path = ?", userID, fileIDList, path).Error
	waitGroup := &sync.WaitGroup{}
	//任务数量
	taskNums := len(userFileIDList)
	waitGroup.Add(taskNums)
	//定义结果集channel
	errChannel := make(chan error, taskNums)

	// 递归删除
	for _, userFileID := range userFileIDList {
		go func(userFileID uint) {
			defer waitGroup.Done()
			// 先看它是否是文件夹
			file := UserFile{}
			err := db.Model(&UserFile{}).Where("id = ? and user_id = ? and file_path = ?", userFileID, userID, path).Scan(&file).Error
			if err != nil {
				errChannel <- err
				return
			}
			//db.Model(&UserFile{}).Delete("file_id = ? and user_id = ? and file_path = ?", fileID, userID, path)
			if file.IsFolder { // 是文件夹
				// 递归进入删除
				//err = DeleteFileDown(userID, []UserFile{file}, path+file.FileName+"/")
				err = DeleteFileDown(userID, file)
				if err != nil {
					errChannel <- err
					return
				}
			}
			// 删除该文件, 不管是否是文件夹
			err = DeleteOneFile(userID, file) // 需要向Recycle表中插入一条记录
			if err != nil {
				errChannel <- err
				return
			}

			errChannel <- err
		}(userFileID)
	}

	//等待所有协程任务完成并关闭结果集通道
	go func() {
		//关闭结果集通道
		defer close(errChannel)
		waitGroup.Wait()
	}()

	// 读取各协程返回的err(查看是否有err)
	for err := range errChannel {
		if err != nil {
			return err
		}
	}

	return nil
}

/*
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
*/

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
	err = db.Unscoped().Model(&file).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	// 更新usedSpace
	var size int64
	db.Model(&File{}).Where("id = ?", file.FileID).Select("size").Scan(&size)
	return UpdateUsedSpace(userID, size)
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
	// deleted_at != null,已经恢复的不用管
	err := db.Unscoped().Model(&UserFile{}).Where("user_id = ? and file_path = ? and deleted_at is not null", userID, path).Scan(&fileList).Error
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
		// 更新usedSpace
		var size int64
		db.Model(&File{}).Where("id = ?", file.FileID).Select("size").Scan(&size)
		err = UpdateUsedSpace(userID, size)
		if err != nil {
			return err
		}
	}
	return nil
}

// 恢复被删除的文件(只有第一层才会将UserFile表中的deleted_at置为NULL)（协程版）
func ResumeFiles(userID uint, userFileIDList []uint) error {
	/**
	file如果是文件夹还要做到：
		1.向上恢复包含它的文件夹
		2.向下递归恢复其子文件
	*/
	waitGroup := &sync.WaitGroup{}
	//任务数量
	taskNums := len(userFileIDList)
	waitGroup.Add(taskNums)
	//定义结果集channel
	errChannel := make(chan error, taskNums)

	for _, userFileID := range userFileIDList {
		go func(userFileID uint) {
			defer waitGroup.Done()
			file := UserFile{}
			db.Unscoped().Model(&UserFile{}).Where("user_id = ? and id = ?", userID, userFileID).Scan(&file)

			// 1.向上恢复包含它的父文件夹(文件和文件夹都要)
			err := ResumeFolderUP(userID, file)
			if err != nil {
				errChannel <- err
				return
			}
			if file.IsFolder { // 是文件夹
				// 2.向下递归恢复其子文件（是由文件夹才要递归下去）
				err = ResumeFilesDown(userID, file)
				if err != nil {
					errChannel <- err
					return
				}
			}
			// 不管是文件夹还是文件，都要将Recycle表中记录删除，(只有第一层)将UserFile表中的deleted_at置为NULL
			err = ResumeOneFile(userID, file)
			if err != nil {
				errChannel <- err
				return
			}
		}(userFileID)
	}

	//等待所有协程任务完成并关闭结果集通道
	go func() {
		//关闭结果集通道
		defer close(errChannel)
		waitGroup.Wait()
	}()

	// 读取各协程返回的err(查看是否有err)
	for err := range errChannel {
		if err != nil {
			fmt.Println("出错啦！！！", err)
			return err
		}
	}

	return nil
}

/*
// 恢复被删除的文件(只有第一层才会将UserFile表中的deleted_at置为NULL)
func ResumeFiles(userID uint, userFileIDList []uint) error {
	//file如果是文件夹还要做到：
	//	1.向上恢复包含它的文件夹
	//	2.向下递归恢复其子文件

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
*/

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

// 移动文件
func MoveFiles(userID uint, fromFileIDList []uint, toFolderPath string) error {
	// 第一层，直接把它们的filePath改为toFolderID的 filePath + fileName + "/"
	//targetFolder := UserFile{}
	//err := db.Model(&UserFile{}).Where("user_id = ? and id = ?", userID, toFolderID).Scan(&targetFolder).Error
	//if err != nil {
	//	return err
	//}
	//newPath := targetFolder.FilePath + targetFolder.FileName + "/"
	for _, fileID := range fromFileIDList {
		file := UserFile{}
		err := db.Model(&UserFile{}).Where("user_id = ? and id = ?", userID, fileID).Scan(&file).Error
		if err != nil {
			return err
		}
		oldPath := file.FilePath + file.FileName + "/"
		err = db.Model(&file).Update("file_path", toFolderPath).Error
		if err != nil {
			return err
		}
		// 递归进入文件夹修改
		if file.IsFolder {
			// 先查出文件夹下的所有文件
			fileIDList := []uint{}
			err = db.Model(&UserFile{}).Where("user_id = ? and file_path = ?", userID, oldPath).Select("id").Scan(&fileIDList).Error
			if err != nil {
				return err
			}
			err = MoveFiles(userID, fileIDList, toFolderPath+file.FileName+"/")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 转存文件
func SaveFiles(userID uint, userFileIDList []uint, savePath string) error {
	for _, userFileID := range userFileIDList {
		file := UserFile{}
		err := db.Model(&UserFile{}).Where("id = ?", userFileID).Scan(&file).Error
		if err != nil {
			return err
		}
		userFile, err := InsertUserFile(userID, file.FileID, file.FileName, savePath, file.IsFolder)
		if err != nil {
			return err
		}

		if file.IsFolder { // 是文件夹，递归
			sonIDList := []uint{}
			err := db.Model(&UserFile{}).Where("file_path = ?", file.FilePath+file.FileName+"/").Select("id").Scan(&sonIDList).Error
			if err != nil {
				return err
			}
			err = SaveFiles(userID, sonIDList, savePath+userFile.FileName+"/")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
// 转存文件
func SaveFiles(userID uint, userFileIDList []uint, savePath string) error {
	for _, userFileID := range userFileIDList {
		file := UserFile{}
		err := db.Model(&UserFile{}).Where("id = ?", userFileID).Scan(&file).Error
		if err != nil {
			return err
		}
		fileName_ := file.FileName
		fileType := ""
		if !file.IsFolder { // 不是文件夹（需要抠出文件后缀）
			list := strings.Split(file.FileName, ".")
			fileName_ = file.FileName[:len(file.FileName)-len(list[len(list)-1])-1]
			fileType = "." + list[len(list)-1]
		}
		i := 0
		for {
			fileName := fileName_
			if i != 0 {
				fileName += fmt.Sprintf("(%v)", i)
			}
			// 查询之前目录下是否已经存在同名文件
			var tmp uint
			db.Model(&UserFile{}).Where("user_id = ? and file_path = ? and file_name = ?", userID, savePath, fileName+fileType).Select("id").First(&tmp)
			if tmp == 0 { // 不存在同名文件
				log.Printf("文件：%v可以插入（没有同名文件）\n", fileName+fileType)
				// 插入新文件
				err = InsertUserFile(userID, file.FileID, fileName+fileType, savePath, file.IsFolder)
				if err != nil {
					return err
				}

				if file.IsFolder { // 是文件夹，递归
					sonIDList := []uint{}
					err := db.Model(&UserFile{}).Where("file_path = ?", file.FilePath+file.FileName+"/").Select("id").Scan(&sonIDList).Error
					if err != nil {
						return err
					}
					err = SaveFiles(userID, sonIDList, savePath+fileName+fileType+"/")
					if err != nil {
						return err
					}
				}

				break
			}
			i++
		}
	}
	return nil
}
*/
