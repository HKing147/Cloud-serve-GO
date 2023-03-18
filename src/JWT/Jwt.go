package JWT

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TokenExpireDuration = time.Hour * 2 // token失效时间

// const TokenExpireDuration = time.Second * 60

var Secret = []byte("网盘系统")

type MyClaims struct {
	UserID uint
	jwt.StandardClaims
}

// get token
func GetToken(userID uint) (string, error) {
	cla := MyClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "HKing",                                    // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	fmt.Println("Token = ", token)
	return token.SignedString(Secret) // 进行签名生成对应的token
}

// parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
