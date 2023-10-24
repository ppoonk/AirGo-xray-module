package router

import (
	"AirGo/api"
	"AirGo/global"
	"AirGo/middleware"
	"AirGo/model"
	"AirGo/web"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var Quit = make(chan os.Signal)

// 初始化总路由
func InitRouter() {
	Router := gin.Default()

	Router.Use(middleware.Serve("/", middleware.EmbedFolder(web.Static, "web"))) // targetPtah=web 是embed和web文件夹的相对路径
	Router.Use(middleware.Cors())
	r := Router.Group("/api")
	{
		r.GET("/ping", func(c *gin.Context) { c.JSON(200, "success") })
	}
	//订阅
	subRouter := r.Group("/sub")
	{
		subRouter.POST("/addSub", api.AddSub)           // 添加订阅
		subRouter.POST("/deleteSub", api.DeleteSub)     // 删除订阅
		subRouter.POST("/getNodeList", api.GetNodeList) // 获取节点列表
		subRouter.POST("/getSubList", api.GetSubList)   // 获取订阅列表
		subRouter.POST("/updateSub", api.UpdateSub)     // 更新订阅
	}
	//节点
	nodeRouter := r.Group("/node")
	{
		nodeRouter.POST("/findNodeById", api.FindNodeById)       // 根据id 查节点
		nodeRouter.POST("/deleteNode", api.DeleteNode)           // 删除节点
		nodeRouter.POST("/updateNode", api.UpdateNode)           // 修改节点
		nodeRouter.POST("/newNode", api.NewNode)                 // 新建节点
		nodeRouter.POST("/tcping", api.Tcping)                   // tcping
		nodeRouter.POST("/getNodePool", api.GetNodePool)         // 获取节点池
		nodeRouter.POST("/joinNodePool", api.JoinNodePool)       // 加入节点池
		nodeRouter.POST("/deleteNodePool", api.DeleteNodePool)   // 从节点池删除节点
		nodeRouter.POST("/getEnabledNodes", api.GetEnabledNodes) // 获取激活的节点
		nodeRouter.POST("/setEnabledNode", api.SetEnabledNode)   // 设置激活的节点
		nodeRouter.GET("/testNodeDelay", api.TestNodeDelay)      // 节点访问网站的延迟
		nodeRouter.GET("/testNodeIP", api.TestNodeIP)            // 节点ip检测

	}
	//shell
	shellRouter := r.Group("/shell")
	{
		shellRouter.GET("/getAllPackages", api.GetAllPackages)
		shellRouter.POST("/doShell", api.DoShell)
		shellRouter.GET("/startService", api.StartService)
		shellRouter.GET("/stopService", api.StopService)
		shellRouter.POST("/getProcessStatus", api.GetProcessStatus)

	}
	//config
	configRouter := r.Group("/shell")
	{
		configRouter.GET("/getConfig", api.GetConfig)
		configRouter.POST("/updateConfig", api.UpdateConfig) // 修改配置
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8899",
		Handler: Router,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logrus.Info("listen: %s\n", err)
			//fmt.Println("listen: %s\n", err)
		}
	}()

	// 等待中断信号关闭服务器(设置 5 秒的超时时间)
	//Quit := make(chan os.Signal)
	Quit <- os.Interrupt
	signal.Notify(Quit, os.Interrupt, os.Kill)
	//signal.Notify(Quit)

	<-Quit
	global.Logrus.Info("Shutdown Server ...")
	//fmt.Println("Shutdown Server ...")
	//退出时关闭xray 清除rules
	global.Logrus.Info("Shutdown xray ...")
	var sh model.Shell
	sh.StopService()
	//
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logrus.Info("Server Shutdown:", err)
		//fmt.Println("Server Shutdown:", err)
	}
	global.Logrus.Info("Server exiting")
	//fmt.Println("Server exiting")

}
