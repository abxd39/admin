package controller

import (
	"admin/app/models/backstage"
	"github.com/gin-gonic/gin"
)

// 权限节点
type NodeController struct {
	BaseController
}

func (n *NodeController) Router(e *gin.Engine) {
	group := e.Group("/node")
	{
		group.GET("/list", n.List)
	}
}

// 节点列表
func (n *NodeController) List(ctx *gin.Context) {
	// 调用model
	modelList, _, err := new(backstage.Node).ListAll(nil)
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
