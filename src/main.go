package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"com/setting"
	"com/gmysql"
	"routers"
	"middlerware"
	"strconv"
	"fmt"
)

func main() {

	var mode = flag.String("mode", "dev", "this is run mode options")
	//命令行解析  mode为*string
	flag.Parse()
	//配置文件初始化
	setting.SetUp(mode)
	//设置模式
	gin.SetMode(setting.ServerSetting.RunMode)
	//数据库初始化
	gmysql.SetUp()

	//创建router对象
	router := gin.New()
	//注册中间件
	router.Use(middlerware.Logger())
	router.Use(gin.Recovery())
	router.Use(middlerware.NotFoundPage())
	//注册路由
	routers.SetUp(router)
	fmt.Println("Listening and serving HTTP on :",setting.ServerSetting.HttpPort)
	fmt.Println("-------------------------------------------")
	router.Run(":" + strconv.Itoa(setting.ServerSetting.HttpPort))
}
