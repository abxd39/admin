package controller

import (
	"admin/app/models/backstage"
	"github.com/gin-gonic/gin"
)

// 权限节点关联API
type NodeAPIController struct {
	BaseController
}

func (n *NodeAPIController) Router(e *gin.Engine) {
	group := e.Group("/node_api")
	{
		group.GET("/list", n.List)
		group.POST("/add", n.Add)
		group.POST("/delete", n.Delete)
	}
}

// 列表
func (n *NodeAPIController) List(ctx *gin.Context) {
	// 获取参数
	nodeId, err := n.GetInt(ctx, "node_id")
	if err != nil || nodeId < 0 {
		n.RespErr(ctx, "参数node_id格式错误")
		return
	}

	// 调用model
	modelList, _, err := new(backstage.NodeAPI).ListAll(nodeId, nil)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	n.Put(ctx, "list", modelList)

	// 返回
	n.RespOK(ctx)
	return
}

// 新增
func (n *NodeAPIController) Add(ctx *gin.Context) {
	// 获取参数
	nodeId, err := n.GetInt(ctx, "node_id")
	if err != nil || nodeId < 0 {
		n.RespErr(ctx, "参数node_id格式错误")
		return
	}

	api := n.GetString(ctx, "api")
	if len(api) == 0 {
		n.RespErr(ctx, "参数api格式错误")
		return
	}

	// 调用model
	nodeAPI := &backstage.NodeAPI{
		NodeId: nodeId,
		Api:    api,
	}
	id, err := new(backstage.NodeAPI).Add(nodeAPI)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	n.Put(ctx, "id", id)

	// 返回
	n.RespOK(ctx)
	return
}

// 删除
func (n *NodeAPIController) Delete(ctx *gin.Context) {
	// 获取参数
	id, err := n.GetInt(ctx, "id")
	if err != nil || id < 0 {
		n.RespErr(ctx, "参数id格式错误")
		return
	}

	// 调用model
	err = new(backstage.NodeAPI).Delete(id)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	// 返回
	n.RespOK(ctx)
	return
}
