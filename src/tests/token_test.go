package tests

import (
	"testing"
	"utils"
	"fmt"
)

var jwtKey = "jwt key"

func TestCreateToken(t *testing.T) {
	userInfo := make(map[string]interface{})
	userInfo["userId"] = 1
	token, _ := utils.CreateToken(jwtKey, userInfo)
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	userInfo := make(map[string]interface{})
	userInfo["userId"] = 1
	token, _ := utils.CreateToken(jwtKey, userInfo)
	fmt.Println("----------parse token--------")
	claims, err := utils.ParseToken(token, "dasda")
	fmt.Println(err)
	fmt.Println(claims)
}
