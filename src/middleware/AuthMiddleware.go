package middleware

import (
	"Cloud-serve/src/JWT"
	"Cloud-serve/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 路由请求中间件，前端必须把token放在请求头上，对服务器进行请求验证token成功后，才能访问后续的请求路由
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("url:", c.Request.URL.Path)
		// 白名单
		WHITE_PATHS := []string{"/api/login", "/api/register", "/api/getCheckCode", "/api/selectUserByID", "/api/getShareByShareUrl", "/api/getFileListByFolderID"}
		exist := false
		for _, path := range WHITE_PATHS {
			if c.Request.URL.Path == path {
				exist = true
				break
			}
		}
		//if c.Request.URL.Path == "/api/login" || c.Request.URL.Path == "/api/register" || c.Request.URL.Path == "/api/getCheckCode" || c.Request.URL.Path == "/api/selectUserByID" || c.Request.URL.Path == "getShareByShareUrl" { // 登录与注册放行
		if exist {
			c.Next()
			return
		}

		// 获取 authorization header：获取前端传过来的信息的
		tokenString_, err := c.Cookie("token")
		if err != nil {
			tokenString_ = ""
		}
		fmt.Println("请求token_:", tokenString_)

		tokenString_admin, err := c.Cookie("admin_token")
		if err != nil {
			tokenString_admin = ""
		}
		fmt.Println("请求token_admin:", tokenString_admin)

		tokenString := tokenString_admin
		if tokenString == "" {
			tokenString = tokenString_
		}
		err = nil
		fmt.Println("请求token:", tokenString)
		// token为空
		if err != nil || tokenString == "" {
			c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "请先登录！！！"}})
			c.Abort()
			return
		}

		// jwt解析token
		claim, err := JWT.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"meta": models.Meta{1, "请先登录！！！"}})
			c.Abort()
			return
		}
		fmt.Println("userID:", claim.UserID)
		c.Set("userID", claim.UserID)
		////验证通过，提取有效部分（除去Bearer)
		//tokenString = tokenString[7:] //截取字符
		////解析token:common/jwt.go
		//token, claims, err := common.ParseToken(tokenString)
		////解析失败||解析后的token无效
		//if err != nil || !token.Valid {
		//	c.JSON(401, gin.H{
		//		"data": gin.H{},
		//		"meta": gin.H{
		//			"msg":  "权限不足",
		//			"code": 401,
		//		},
		//	})
		//	return
		//}
		//
		////token通过验证, 获取claims中的UserID
		//userId := claims.UserId
		//var user model.User
		////查询数据库
		//common.DB.First(&user, userId)
		//
		//// 验证用户是否存在
		//if user.ID == 0 {
		//	c.JSON(401, gin.H{
		//		"data": gin.H{},
		//		"meta": gin.H{
		//			"msg":  "权限不足",
		//			"code": 401,
		//		},
		//	})
		//}
		//
		////用户存在 将user信息写入上下文
		//c.Set("user", user)

		c.Next()
	}
}
