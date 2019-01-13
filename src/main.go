package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"com/setting"
	"com/gmysql"
)

func main() {

	var mode = flag.String("mode", "dev", "this is run mode options")
	//命令行解析  mode为*string
	flag.Parse()
	//配置文件初始化
	setting.SetUp(mode)
	//数据库初始化
	gmysql.SetUp()
	//创建router对象
	router := gin.New()
	fmt.Println(router)
}
