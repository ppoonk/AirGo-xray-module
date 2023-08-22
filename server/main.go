package main

import (
	"AirGo/global"
	"AirGo/initialize"
	"AirGo/model"
	"AirGo/router"
	"flag"
	"fmt"
	"runtime"
)

var start = flag.Bool("start", false, "启动")
var stop = flag.Bool("stop", false, "停止")

func main() {
	switch runtime.GOOS {
	case "darwin":
		initialize.Initialize() //初始化全局变量，配置参数
		router.InitRouter()     //初始化路由
	default:
		flag.Parse()
		if *start && !*stop { //启动
			initialize.Initialize() //初始化参数

			var sh model.Shell
			sh.DoShell(model.OpenFirewall, false) //linux放行防火墙

			router.InitRouter() //初始化路由
			//判断是否自启xray
			global.Logrus.Info("判断是否自启xray:", global.Config.StartupXray)
			if global.Config.StartupXray == "1" { //是否开机自启xray
				err := sh.StartService()
				if err != nil {
					global.Logrus.Error("自启xray error:", err.Error())
				}
			}

		} else if !*start && *stop { //停止
			var sh model.Shell
			sh.StopService()
			fmt.Println("停止")
		} else {
			fmt.Println("非法参数，退出")
		}
	}
}
