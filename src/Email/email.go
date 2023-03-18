package Email

import (
	"Cloud-serve/src/DB"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

// T: 验证码有效期
func SendEmail(emailTo string, T int) (err error) {
	config := map[string]string{
		"user": "1470042308@qq.com", // 邮件发送者的地址
		"pass": "rmtxhlaukqhthigf",  // qq邮箱填授权码，百度一下获取方式。
		"host": "smtp.qq.com",       // 发送将邮件发送给腾讯的smtp邮件服务器
	}
	// 生成6位数字验证码
	checkCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	body := fmt.Sprintf("您的验证码为：【%v】，请勿泄露，祝您生活愉快！", checkCode)
	//rd := DB.ConnRedis()
	//ctx := context.Background()
	//rd.Set(ctx, emailTo, checkCode, time.Duration(T)*time.Second) // checkCode写入redis,T秒后过期
	DB.Set(emailTo, checkCode, time.Duration(T)*time.Second) // checkCode写入redis,T秒后过期

	m := gomail.NewMessage()
	m.SetHeader("From", config["user"], "【阿里云盘】")
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", "来自【阿里云盘】的验证码") //邮件标题
	m.SetBody("text/html", body)           //设置邮件正文

	d := gomail.NewDialer(config["host"], 465, config["user"], config["pass"])
	err = d.DialAndSend(m) // 发送邮件

	return err
}
