package main

import (
	"flag"
	"fmt"
	"runtime"
	"server/global"
	"server/initialize"
	"server/model"
	"server/router"
)

var start = flag.Bool("start", false, "启动")
var stop = flag.Bool("stop", false, "停止")

func main() {
	switch runtime.GOOS {
	case "darwin":
		initialize.Initialize() //初始化参数
		router.InitRouter()     //初始化路由
	default:
		flag.Parse()
		if *start && !*stop { //启动
			initialize.Initialize() //初始化参数
			if global.Config.OS != "darwin" {
				var sh model.Shell
				sh.DoShell(model.OpenFirewall, false)
			} //放行防火墙
			router.InitRouter() //初始化路由

		} else if !*start && *stop { //停止
			var sh model.Shell
			sh.StopService()
			fmt.Println("停止")
		} else {
			fmt.Println("非法参数，退出")
		}
	}
}
