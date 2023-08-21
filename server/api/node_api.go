package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/model/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 根据id 查节点
func FindNodeById(ctx *gin.Context) {
	nodeIdStr := ctx.Query("nodeId")
	if nodeIdStr == "" {
		return
	}
	nodeId, _ := strconv.ParseInt(nodeIdStr, 10, 64)
	var node = model.Node{ID: nodeId}
	err := node.FindNodeById()
	if err != nil {
		global.Logrus.Error("查节点错误:", err.Error())
		response.Fail("查节点错误:"+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("node:", node)
	response.OK("查节点成功", node, ctx)
}

// 删除节点
func DeleteNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("删除节点参数错误:", err.Error())
		response.Fail("删除节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = node.DeleteNode()
	if err != nil {
		global.Logrus.Error("删除节点错误:", err.Error())
		response.Fail("删除节点错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("删除节点成功", nil, ctx)
}

// 修改节点
func UpdateNode(ctx *gin.Context) {
	var node model.NodePool
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("修改节点参数错误:", err.Error())
		response.Fail("修改节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = node.UpdateNode()
	if err != nil {
		global.Logrus.Error("修改节点错误:", err.Error())
		response.Fail("修改节点错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改节点成功", nil, ctx)
}

// 新建节点
func NewNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("新建节点参数错误:", err.Error())
		response.Fail("新建节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = node.NewNode()
	if err != nil {
		global.Logrus.Error("新建节点错误:", err.Error())
		response.Fail("新建节点错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("新建节点成功", nil, ctx)
}

// 节点tcping
func Tcping(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("节点tcping参数错误:", err.Error())
		response.Fail("节点tcping参数错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("节点tcping成功", node.Tcping(), ctx)
}

// 获取节点池
func GetNodePool(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("获取节点池参数错误:", err.Error())
		response.Fail("获取节点池参数错误:"+err.Error(), nil, ctx)
		return
	}
	res, err := node.GetNodePool()
	if err != nil {
		global.Logrus.Error("获取节点池错误:", err.Error())
		//response.Fail("获取节点池错误:"+err.Error(), nil, ctx)
		//return
	}
	response.OK("获取节点池成功", res, ctx)
}

// 加入节点池
func JoinNodePool(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("加入节点池参数错误:", err.Error())
		response.Fail("加入节点池参数错误:"+err.Error(), nil, ctx)
		return
	}
	fmt.Println("加入节点池:", node)
	err = node.JoinNodePool()
	if err != nil {
		global.Logrus.Error("加入节点池错误:", err.Error())
		response.Fail("加入节点池错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("加入节点池成功", nil, ctx)
}

// 从节点池删除节点
func DeleteNodePool(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("从节点池删除节点参数错误:", err.Error())
		response.Fail("从节点池删除节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = node.DeleteNodePool()
	if err != nil {
		global.Logrus.Error("从节点池删除节点错误:", err.Error())
		response.Fail("从节点池删除节点错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("从节点池删除节点成功", nil, ctx)
}

// 获取激活的节点
func GetEnabledNodes(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("获取激活的节点参数错误:", err.Error())
		response.Fail("获取激活的节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	nodes, err := node.GetEnabledNodes()
	if err != nil {
		//if err == gorm.ErrRecordNotFound {
		//	global.Logrus.Error("未设置活动节点:", err.Error())
		//	response.Fail("未设置活动节点", nil, ctx)
		//	return
		//} else {
		//	global.Logrus.Error("获取激活的节点错误:", err.Error())
		//	response.Fail("错误:"+err.Error(), nil, ctx)
		//	return
		//}
		global.Logrus.Error("获取激活的节点，error:", err.Error())
	}
	response.OK("获取激活的节点成功", nodes, ctx)
}

// 设置激活节点
func SetEnabledNode(ctx *gin.Context) {
	var node model.Node
	err := ctx.ShouldBind(&node)
	if err != nil {
		global.Logrus.Error("设置激活节点参数错误:", err.Error())
		response.Fail("设置激活节点参数错误:"+err.Error(), nil, ctx)
		return
	}
	err = node.SetEnableNodeToDb()
	if err != nil {
		global.Logrus.Error("设置激活节点错误:", err.Error())
		response.Fail("设置激活节点错误:"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取激活节点成功", nil, ctx)

}

// 节点访问网站的延迟
func TestNodeDelay(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		response.Fail("url为空", nil, ctx)
		return
	}
	var node model.Node
	response.OK("节点访问网站的延迟", node.TestNodeDelayNew(url), ctx)
}

// 节点ip检测
func TestNodeIP(ctx *gin.Context) {
	t := ctx.Query("type")
	if t == "" {
		response.Fail("节点ip检测type为空", nil, ctx)
		return
	}
	var node model.Node

	res, err := node.TestNodeIP(t)
	if err != nil {
		global.Logrus.Error("节点ip检测错误:", err.Error())
		response.Fail("节点ip检测错误:"+err.Error(), nil, ctx)
		return
	}

	response.OK("节点ip检测", res, ctx)

}
