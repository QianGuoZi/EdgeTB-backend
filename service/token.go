package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//密钥
var jwtSecret = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 根据用户的用户名产生token
func GenerateToken(username string) (string, time.Time, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * time.Hour)

	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "EdgeTB",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, expireTime, err
}

// ParseToken 根据传入的token值获取到Claims对象信息，（进而获取其中的信息）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
	// 要传入指针，项目中结构体都是用指针传递，节省空间。
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
