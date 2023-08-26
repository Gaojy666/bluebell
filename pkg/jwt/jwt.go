package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("夏天夏天悄悄过去")

type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己声明的数据
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bluebell",                                 // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的claims则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid {
		// 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
