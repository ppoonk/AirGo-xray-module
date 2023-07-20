package api

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/model"
	"server/model/response"
)

// 获取全部包名
func GetAllPackages(ctx *gin.Context) {
	var sh model.Shell
	out, err := sh.GetAllPackages()
	if err != nil {
		global.Logrus.Error("获取全部包名错误:", err.Error())
		response.Fail("获取全部包名错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取全部包名成功", out, ctx)
}

// 执行shell
func DoShell(ctx *gin.Context) {
	//r, _ := io.ReadAll(ctx.Request.Body)
	//fmt.Println("请求：", string(r))
	var sh model.Shell
	err := ctx.ShouldBind(&sh)
	if err != nil {
		global.Logrus.Error("执行shell错误:", err.Error())
		response.Fail("执行shell错误:"+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("shell：", sh)
	out, err := sh.DoShell(sh.Shell, sh.OutType)
	if err != nil {
		global.Logrus.Error("执行shell错误:", err.Error())
		response.Fail("执行shell错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("执行shell成功", out, ctx)
}

// 启动服务
func StartService(ctx *gin.Context) {
	var sh model.Shell
	err := sh.StartService()
	if err != nil {
		global.Logrus.Error("启动服务错误:", err.Error())
		response.Fail("启动服务错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("启动服务成功", nil, ctx)

}

// 关闭服务
func StopService(ctx *gin.Context) {
	var sh model.Shell
	out, err := sh.StopService()
	if err != nil {
		global.Logrus.Error("关闭服务错误:", err.Error())
		response.Fail("关闭服务错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("关闭服务", out, ctx)
}

// 查询进程状态
func GetProcessStatus(ctx *gin.Context) {
	var sh model.Shell
	err := ctx.ShouldBind(&sh)
	if err != nil {
		global.Logrus.Error("查询进程状态错误:", err.Error())
		response.Fail("查询进程状态错误:"+err.Error(), nil, ctx)
		return
	}
	out, err := sh.GetProcessStatus()
	global.Logrus.Info("查询进程状态")
	if err != nil {
		if err.Error() == "exit status 1" {
			response.Fail(sh.Shell+"未启动:"+err.Error(), nil, ctx)
			return
		} else {
			global.Logrus.Error("查询进程状态错误:", err.Error())
			response.Fail("查询进程状态错误:"+err.Error(), nil, ctx)
			return
		}
	}
	if len(out) > 2 {
		response.OK("查询进程状态成功", out, ctx)
		return

	} else {
		response.Fail(sh.Shell+"未启动", nil, ctx)
		return
	}

}
