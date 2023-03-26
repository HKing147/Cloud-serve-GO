package models

import (
	"Cloud-serve/src/DB"
	"Cloud-serve/src/JWT"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// User
type User struct {
	gorm.Model
	Avatar string `json:"avatar" gorm:"avatar"` // 头像
	Email  string `json:"email" gorm:"email"`   // 邮箱
	//ID           int64  `json:"id" gorm:"id"`                      // 用户ID
	Password string `json:"password" gorm:"password"` // 密码
	QQ       string `json:"QQ" gorm:"qq"`             // QQ号
	//RegisterTime int64  `json:"registerTime" gorm:"register_time"` // 注册时间
	TotalSpace int64  `json:"totalSpace" gorm:"total_space;default:10737418240"` // 网盘容量(10G:10*1024*1024*1024)
	UsedSpace  int64  `json:"usedSpace" gorm:"used_space"`                       // 已用容量
	UserName   string `json:"userName" gorm:"user_name"`                         // 用户名
	Wechat     string `json:"Wechat" gorm:"wechat"`                              // 微信号
}

type UserRegisterForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CheckCode string `json:"checkCode"`
}

type UserLoginForm struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

func Register(c *gin.Context) {
	userRegisterForm := UserRegisterForm{}
	c.BindJSON(&userRegisterForm)
	fmt.Println("userRegisterForm:", userRegisterForm)
	// 查看邮箱是否已注册
	if _, err := SelectUserByEmail(userRegisterForm.Email); err == nil { // 邮箱已注册
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "邮箱已注册！！！"}})
		return
	}
	// 检查验证码是否正确
	rdCheckCode, _ := DB.Get(userRegisterForm.Email).Result()
	if rdCheckCode != userRegisterForm.CheckCode || len(rdCheckCode) != 6 {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "验证码错误！！！"}})
		return
	}
	// 邮箱未注册
	user := &User{}
	user.UserName = userRegisterForm.Email
	user.Email = userRegisterForm.Email
	// 加密密码
	h := md5.New()
	h.Write([]byte(userRegisterForm.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))
	err := InsertUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "注册失败"}})
		return
	}
	// 生成token
	token, err := JWT.GetToken(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "注册失败"}})
		return
	}
	// Token写入Redis
	DB.Set(user.Email+"_token", token, time.Hour) // 有效期一小时
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", token, 60*60, "/", "http://localhost:5173", true, false)
	c.JSON(http.StatusOK, gin.H{"meta": Meta{0, "注册成功"}, "token": token})
}

func Login(c *gin.Context) {
	userLoginForm := UserLoginForm{}
	c.BindJSON(&userLoginForm)
	fmt.Println("userLoginForm:", userLoginForm)
	// 查看email是否已注册
	user, err := SelectUserByEmail(userLoginForm.Email)
	if err != nil { // 邮箱未注册
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "登录失败，邮箱或密码错误！！！"}})
		return
	}
	// 邮箱已注册
	// 查看密码是否正确
	// 加密密码
	h := md5.New()
	h.Write([]byte(userLoginForm.Password))
	if hex.EncodeToString(h.Sum(nil)) != user.Password { // 密码错误
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "登录失败，邮箱或密码错误！！！"}})
		return
	}
	// 邮箱和密码正确
	// 生成token
	token, err := JWT.GetToken(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"meta": Meta{1, "登录失败！！！"}})
		return
	}
	// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjMsImV4cCI6MTY3ODk1NDQ3MywiaXNzIjoiSEtpbmcifQ.2k5H2vvG_QzkN6KFrfja6M32tpLxUKWflauDmmJ1VvY"
	// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjMsImV4cCI6MTY3ODk1NDkwMCwiaXNzIjoiSEtpbmcifQ.XKUA1OsR4z_ftSZy1RjlgXpvYb2uSyT8Kzl_z6Mnyj4"
	// Token写入Redis
	DB.Set(user.Email+"_token", token, time.Hour) // 有效期一小时
	//c.SetSameSite(http.SameSiteNoneMode)
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", token, 60*60, "/", "http://localhost:5173", true, false)
	c.JSON(http.StatusOK, gin.H{"meta": Meta{0, "登录成功！"}, "token": token})
}

func InsertUser(user *User) error {
	res := db.Create(user)
	return res.Error
}

func SelectUserByEmail(email string) (user User, err error) {
	result := db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func SelectUserByUserName(userName string) (user User, err error) {
	result := db.Where("user_name = ?", userName).First(&user)
	return user, result.Error
}

func SelectUserByID(userID uint) (User, error) {
	user := User{}
	result := db.Where("id = ?", userID).First(&user)
	return user, result.Error
}

// 使用容量更新(+/-)
func UpdateUsedSpace(userID uint, delta int64) error {
	return db.Model(&User{}).Where("id = ?", userID).Update("used_space", gorm.Expr("used_space + ?", delta)).Error
}
