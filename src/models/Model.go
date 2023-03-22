package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(
		mysql.Open("root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)}) //配置日志级别，打印出所有的sql

	printError("连接数据库失败", err)
	err = db.AutoMigrate(&User{}, &File{}, &UserFile{}, &Recycle{}, &Share{})
	printError("建表失败", err)
}

func printError(str string, err error) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}
