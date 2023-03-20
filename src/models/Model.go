package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(
		mysql.Open("root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{})
	printError("连接数据库失败", err)
	err = db.AutoMigrate(&User{}, &File{}, &UserFile{}, &Recycle{})
	printError("建表失败", err)
}

func printError(str string, err error) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}
