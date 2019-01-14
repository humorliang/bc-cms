package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"com/setting"
)

//创建一个jwt token
func CreateToken(key string, mUser map[string]interface{}) (string, error) {
	//指定加密方式产生一个token对象
	token := jwt.New(jwt.SigningMethodHS256)
	//创建一个claims
	claims := make(jwt.MapClaims)

	//标准claims添加
	//token生效时间
	now := time.Now()
	expireTime := now.Add(setting.ServerSetting.TokenTimeoutHour)
	claims["iss"] = "humorliang"      //jwt签发者
	claims["sub"] = "all"             //jwt所面向的用户
	claims["aud"] = "web"             //接收jwt的一方
	claims["exp"] = expireTime.Unix() //jwt的过期时间，这个过期时间必须要大于签发时间 使用时间戳
	//claims["nbf"] = ""                //定义在什么时间之前，该jwt都是不可用的.
	claims["iat"] = now.Unix()        //jwt的签发时间
	claims["jti"] = now.Unix()        //jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击

	//私有claims添加
	for index, val := range mUser {
		claims[index] = val
	}
	//设置token的claims
	token.Claims = claims
	//返回一个token
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//解析一个jwt token
//返回一个claims error
func ParseToken(tokenString string, key string) (map[string]interface{}, error) {
	//解析得到一个token
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
			}
			return []byte(key), nil
		})
	if err != nil {
		return nil, err
	}
	//对token的claims进行解析，token验证
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token is parse and valid fails! ")
	}
}
