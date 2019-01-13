package file

import (
	"path"
	"os"
	"fmt"
)

//获取文件扩展名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

//判断资源是否存在
func CheckNotExist(src string) bool {
	//返回文件信息 err 会返回*PathError
	_, err := os.Stat(src)
	//返回一个布尔值 说明 该错误 是否 表示一个文件或目录不存在
	return os.IsNotExist(err)
}

//创建文件夹
func MkDir(src string) error {
	//创建文件夹组 os.Mkdir只会创建一个文件夹
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//判断文件夹并创建
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

//文件打开创建追加
func OpenCreateAppend(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//检查权限
func CheckPermission(src string) bool {
	//返回资源描述信息
	_, err := os.Stat(src)
	//异常没权限 true
	return os.IsPermission(err)
}

//指定文件夹创建文件，不存在则创建
func MustOpenSrc(fileName string, fileDir string) (*os.File, error) {
	//获取当前程序runtime的路径
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.GetWd error:%s", err)
	}
	//获取文件夹路径
	src := dir + "/" + fileDir
	//检验权限
	if CheckPermission(src) {
		return nil, fmt.Errorf("permission is denied src:%s", src)
	}
	//检验文件夹并创建
	if err := IsNotExistMkDir(src); err != nil {
		return nil, fmt.Errorf("mkdir is error:%s", err)
	}
	//打开文件流
	f, err := OpenCreateAppend(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file error:%s", err)
	}
	return f, nil
}
