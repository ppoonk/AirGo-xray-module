package utils

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// 获取配置文件所在目录
func GetConfigDir() string {
	dir := os.Getenv("AirGo")
	if dir != "" && IsDir(dir) {
		return dir
	}
	return GetRunPath()
}

// 获取当前进程的可执行文件的路径名
func GetRunPath() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}

// 判断路径是否正确且为文件夹
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return i.IsDir()
	}
	return false
}

// WriteJSON 将对象写入json文件
func WriteJSON(v interface{}, path string) error {
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	return jsonEncode.Encode(v)
}

// 设置日志文件
func SetOutputFile() (*os.File, error) {
	now := time.Now()
	logFileName := now.Format("2006-01-02") + ".log" //日志文件名
	logFilePath := ""                                //路径
	if dir, err := os.Getwd(); err == nil {          //当前工作目录
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) { //isNotExist()判断为true，说明文件或者文件夹不存在
		if err := os.MkdirAll(logFilePath, 0777); err != nil { //创建
			log.Println(err.Error())
			return nil, err
		}
	}

	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //如果已经存在，则在尾部添加写
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
