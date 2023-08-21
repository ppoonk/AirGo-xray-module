package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/model/response"
	"github.com/gin-gonic/gin"
)

// 获取配置
func GetConfig(ctx *gin.Context) {
	var c model.Config
	err := c.GetConfig()
	if err != nil {
		global.Logrus.Error("获取配置错误:", err.Error())
		response.Fail("获取配置错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取配置成功", c, ctx)
}

// 修改配置
func UpdateConfig(ctx *gin.Context) {
	var c model.Config
	err := ctx.ShouldBind(&c)
	if err != nil {
		global.Logrus.Error("修改配置错误:", err.Error())
		response.Fail("修改配置错误:"+err.Error(), nil, ctx)
		return
	}
	err = c.UpdateConfig()
	if err != nil {
		global.Logrus.Error("修改配置错误:", err.Error())
		response.Fail("修改配置错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改配置成功", nil, ctx)
}
