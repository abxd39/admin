package controller

import (
	"admin/app/models/backstage"
	"github.com/gin-gonic/gin"
)

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
func (n *NodeController) List(c *gin.Context) {
	// 获取参数
	page, err := n.GetInt(c, "page", 1)
	if err != nil {
		n.RespErr(c, "参数page格式错误")
		return
	}

	rows, err := n.GetInt(c, "rows", 10)
	if err != nil {
		n.RespErr(c, "参数rows格式错误")
		return
	}

	// 调用model
	list, err := new(backstage.Node).ListAll(page, rows, nil)
	if err != nil {
		n.RespErr(c, err)
		return
	}

	// 设置返回数据
	n.Put(c, "list", list)

	// 返回
	n.RespOK(c)
	return
}
