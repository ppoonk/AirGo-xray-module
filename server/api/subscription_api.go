package api

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/model"
	"server/model/response"
)

// 添加订阅
func AddSub(ctx *gin.Context) {
	var sub model.Subscription
	err := ctx.ShouldBind(&sub)
	if err != nil {
		global.Logrus.Error("添加订阅参数错误:", err.Error())
		response.Fail("添加订阅参数错误:", err.Error(), ctx)
		return
	}
	//fmt.Println("添加订阅sub.Url:", sub.Url)
	err = sub.ParseSub()
	if err != nil {
		global.Logrus.Error("添加订阅错误:", err.Error())
		response.Fail("添加订阅错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("添加订阅成功", nil, ctx)
}

// 更新订阅
func UpdateSub(ctx *gin.Context) {
	var sub model.Subscription
	err := ctx.ShouldBind(&sub)
	if err != nil || sub.ID == 1 {
		global.Logrus.Error("更新订阅参数错误:", err.Error())
		response.Fail("更新订阅参数错误:", err.Error(), ctx)
		return
	}
	err = sub.UpdateSub()
	if err != nil {
		global.Logrus.Error("更新订阅错误:", err.Error())
		response.Fail("更新订阅错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("更新订阅成功", nil, ctx)
}

// 删除订阅
func DeleteSub(ctx *gin.Context) {
	var sub model.Subscription
	err := ctx.ShouldBind(&sub)
	if err != nil {
		global.Logrus.Error("删除订阅参数错误:", err.Error())
		response.Fail("删除订阅参数错误:", err.Error(), ctx)
		return
	}
	err = sub.DeleteSub()
	if err != nil {
		global.Logrus.Error("删除订阅错误:", err.Error())
		response.Fail("删除订阅错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("删除订阅成功", nil, ctx)
}

// 获取节点列表
func GetNodeList(ctx *gin.Context) {
	var sub model.Subscription
	err := ctx.ShouldBind(&sub)
	if err != nil {
		global.Logrus.Error("获取节点列表参数错误:", err.Error())
		response.Fail("获取节点列表参数错误:", err.Error(), ctx)
		return
	}
	err = sub.GetSubNodeList()
	if err != nil {
		global.Logrus.Error("获取节点列表错误:", err.Error())
		response.Fail("获取节点列表错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取节点列表成功", sub.Nodes, ctx)

}

// 获取订阅列表
func GetSubList(ctx *gin.Context) {
	var sub model.Subscription
	err := ctx.ShouldBind(&sub)
	if err != nil {
		global.Logrus.Error("获取订阅列表参数错误:", err.Error())
		response.Fail("获取订阅列表参数错误:", err.Error(), ctx)
		return
	}
	res, err := sub.GetSubList()
	if err != nil {
		global.Logrus.Error("获取订阅列表错误:", err.Error())
		response.Fail("获取订阅列表错误:", err.Error(), ctx)
		return
	}
	response.OK("获取订阅列表成功", res, ctx)
}
