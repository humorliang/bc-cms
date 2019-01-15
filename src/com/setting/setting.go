package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

//app结构体
type App struct {
	RuntimePath string
	LogPath     string
	JwtKey      string
}

//server结构体
type Server struct {
	RunMode          string
	HttpPort         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	TokenTimeoutHour time.Duration
}

//database结构体
type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var (
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DatabaseSetting = &Database{}
	//定义一个配置文件流
	cfg *ini.File
)

//配置文件初始化函数
func SetUp(mode *string) {
	var err error
	// 相应环境读取配置文件
	if *mode == "pro" {
		cfg, err = ini.Load("conf/pro.ini")
		if err != nil {
			log.Fatalf("setting load file ‘conf/pro.ini’ error:%s", err)
		}
	} else {
		cfg, err = ini.Load("conf/dev.ini")
		if err != nil {
			log.Fatalf("setting load file ‘conf/dev.ini’ error:%s", err)
		}
	}
	sectionToMap("app", AppSetting)
	sectionToMap("server", ServerSetting)
	sectionToMap("database", DatabaseSetting)
	//相对应的配置
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	ServerSetting.TokenTimeoutHour = ServerSetting.TokenTimeoutHour * time.Hour
}

//文件信息与结构体绑定
func sectionToMap(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("setting section to map error:%s", err)
	}
}
