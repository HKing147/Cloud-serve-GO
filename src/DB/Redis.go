package DB

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type User struct {
	Name    string
	Age     int
	Address string
}

func (user *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(user)
}

func (user *User) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, user)
}

//user := User{
//Name:    "小明",
//Age:     18,
//Address: "北京",
//}
//rd.Set("user1", &user, 10*time.Second) // 10秒过期，为0表示永不过期
//user2 := User{}
//rd.Get("user1").Scan(&user2)

func ConnRedis() *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr: "47.100.190.87:6379", // url
		//Password: "123456",
		DB: 0, // 0号数据库
	})
	result, err := rd.Ping().Result()
	if err != nil {
		fmt.Println("ping err :", err)
		return nil
	}
	fmt.Println(result)
	return rd
}

func Test() {
	rd := ConnRedis()

	user1 := User{
		Name:    "小明",
		Age:     18,
		Address: "北京",
	}
	user2 := User{
		Name:    "Tom",
		Age:     18,
		Address: "北京",
	}

	rd.SAdd("users", &user1)
	rd.SAdd("users", &user2)
	rd.Expire("users", 10*time.Second)
	users := []User{}
	rd.SMembers("users").ScanSlice(&users)
	for _, user := range users {
		fmt.Println(user.Name, user.Age, user.Address)
	}
}
