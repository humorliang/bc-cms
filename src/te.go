package main

import (
	"com/gmysql"
	"fmt"
	"db"
)

func main() {
	//t := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(t)
	////返回当前程序运行的目录的根路径
	//src,_:=os.Getwd()
	//fmt.Println(strings.Split(src,"/")[:len(strings.Split(src,"/"))-1])
	////fmt.Errorf("cuowu")
	//logging.Debug("debug test")
	//jwtKey := "test key"
	//type user map[string]interface{}
	//userInfo := make(user)
	//userInfo["userId"] = 1
	//token, err := utils.CreateToken(jwtKey, userInfo)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(token)
	//claims, err := utils.ParseToken(token, jwtKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(claims)
	//a:=[]interface{}{1,"2",&jwtKey}
	//fmt.Println(a)
	rows, err := gmysql.Con.Query("select user_id,user_login from bc_user where user_login='u1' ")
	if err != nil {
		fmt.Println(err)
	}
	rs, err := db.Querys(rows)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(rs)
}
