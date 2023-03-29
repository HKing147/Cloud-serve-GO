package DB

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//type User struct {
//	Name    string
//	Age     int
//	Address string
//}
//
//func (user *User) MarshalBinary() (data []byte, err error) {
//	return json.Marshal(user)
//}
//
//func (user *User) UnmarshalBinary(data []byte) (err error) {
//	return json.Unmarshal(data, user)
//}

// user := User{
// Name:    "小明",
// Age:     18,
// Address: "北京",
// }
// rd.Set("user1", &user, 10*time.Second) // 10秒过期，为0表示永不过期
// user2 := User{}
// rd.Get("user1").Scan(&user2)

//	func ConnRedis() *redis.Client {
//		rd := redis.NewClient(&redis.Options{
//			Addr:     "47.100.190.87:6379", // url
//			Password: "flzx3qc",
//			DB:       0, // 0号数据库
//		})
//		ctx := context.Background()
//		result, err := rd.Ping(ctx).Result()
//		if err != nil {
//			fmt.Println("ping err :", err)
//			return nil
//		}
//		fmt.Println(result)
//		return rd
//	}
var rd *redis.Client
var ctx context.Context

func InitRedis() {
	rd = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 发布
		//Addr:     "47.100.190.87:6379", // 开发
		Password: "flzx3qc",
		DB:       0, // 0号数据库
	})
	ctx = context.Background()
	result, err := rd.Ping(ctx).Result()
	if err != nil {
		fmt.Println("ping err :", err)
		return
	}
	fmt.Println(result)
}

func Get(key string) *redis.StringCmd {
	return rd.Get(ctx, key)
}

func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return rd.Set(ctx, key, value, expiration)
}

func SMembers(key string) *redis.StringSliceCmd {
	return rd.SMembers(ctx, key)
}

func Exists(key ...string) *redis.IntCmd {
	return rd.Exists(ctx, key...)
}

func SAdd(key string, members ...interface{}) *redis.IntCmd {
	return rd.SAdd(ctx, key, members)
}

func Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return rd.Expire(ctx, key, expiration)
}

//func Test() {
//	rd := ConnRedis()
//
//	user1 := User{
//		Name:    "小明",
//		Age:     18,
//		Address: "北京",
//	}
//	user2 := User{
//		Name:    "Tom",
//		Age:     18,
//		Address: "北京",
//	}
//
//	rd.SAdd("users", &user1)
//	rd.SAdd("users", &user2)
//	rd.Expire("users", 10*time.Second)
//	users := []User{}
//	rd.SMembers("users").ScanSlice(&users)
//	for _, user := range users {
//		fmt.Println(user.Name, user.Age, user.Address)
//	}
//}
