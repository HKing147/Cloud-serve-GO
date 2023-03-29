# Cloud-serve-GO
基于对象存储的网盘系统后端-Golang

# 打包部署前的修改
1. `Redis`地址修改: [Redis.go](./src/DB/Redis.go)中修改`Addr`
2. `Mysql`地址修改: [Model.go](./src/models/Model.go)中修改`Addr
3. [CORSMiddleware.go](./src/middleware/CORSMiddleware.go)的`Access-Control-Allow-Origin`改为`localhost`
4. [User.go](./src/models/User.go)的`Register`与`Login`方法中的`c.SetSameSite(http.SameSiteNoneMode)`注释掉，`c.SetCookie`改为服务器的。
