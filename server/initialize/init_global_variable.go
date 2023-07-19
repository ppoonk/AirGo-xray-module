package initialize

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"server/global"
	"server/model"
	"server/utils"
)

func Initialize() {
	global.DB = InitGormSqlite() //连接数据库
	global.Logrus = InitLogrus() //初始化logrus
	if global.DB != nil {        //初始化数据库
		RegisterTables()
		InsertInto()
	}
	NodeCrontab() //节点连通性检测定时任务
	//InitClientAndDialer() //初始化http.Client,net.Dialer
	InitConfig() //初始化配置
}

//func InitClientAndDialer() {
//	global.ClientWithDNS = utils.ClientWithDNS("114.114.114.114", 5*time.Second)
//	global.ClientWithSocks5 = utils.ClientWithSocks5("localhost", 1231, 5*time.Second)
//	global.Dialer = utils.Dialer("114.114.114.114", 5*time.Second)
//}

func InitLogrus() *logrus.Logger {
	//实例化
	logger := logrus.New()
	src, _ := utils.SetOutputFile()
	//设置输出(同时输出到标准输出和日志文件)
	mw := io.MultiWriter(os.Stdout, src)
	logger.SetOutput(mw)
	//设置输出
	//logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

// 节点连通性检测定时任务
func NodeCrontab() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("*/30 * * * * *", func() {
		if global.Config.AutoChangeNode == "0" {
			//global.Logrus.Info("AutoChangeNode=", global.Config.AutoChangeNode, "跳过切换节点")
			return
		}
		var sh model.Shell
		out, err := sh.GetProcessStatus()
		if len(out) < 3 || err != nil {
			//xray 关闭，跳过
			//global.Logrus.Info("xray关闭，跳过切换节点")
			return
		}
		//global.Logrus.Info("节点连通性检测定时任务开始")
		var node1 = model.Node{Ascription: "domestic"}
		var node2 = model.Node{Ascription: "abroad"}
		ok1 := node1.SetEnableNode()
		ok2 := node2.SetEnableNode()
		if ok1 || ok2 {
			//重启xray
			sh.StopService()
			sh.StartService()
		}
	})
	if err != nil {
		return
	}
	global.Logrus.Info("节点连通性检测定时任务")
	c.Start()
}

func InitConfig() {
	var c model.Config
	err := c.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("global.Config", global.Config)
}
