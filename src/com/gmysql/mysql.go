package gmysql

import (
	"database/sql"
	"fmt"
	"com/logging"
	"com/setting"
)

var Con *sql.DB

func SetUp() {
	var err error
	//数据库连接地址
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)
	//打开一个连接池
	Con, err = sql.Open(setting.DatabaseSetting.Type, dataSource)
	if err != nil {
		logging.Fatal("数据库连接错误：%s", err)
	}
	//测试数据库
	err = Con.Ping()
	if err != nil {
		logging.Fatal("数据库测试连接失败：%s", err)
	}
}
